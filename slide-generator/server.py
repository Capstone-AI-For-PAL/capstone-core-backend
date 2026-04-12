from flask import Flask, request, send_file
import subprocess
import os
import uuid
import traceback
import logging
import sys
import io

from helper import validateSectionInput
from pptx_generator import generate_pptx

app = Flask(__name__)

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s - %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger("slide-generator")

SUPPORTED_FILE_TYPES = {"pdf", "pptx"}
MARP_FLAG = {"pdf": "--pdf", "pptx": "--pptx-editable"}
MIME_TYPE = {
    "pdf": "application/pdf",
    "pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
}


@app.route('/generate', methods=['POST'])
def generate_slides():
    input_filename = None
    output_filename = None

    try:
        data = request.json
        if not data:
            return {"error": "Invalid JSON body"}, 400

        markdown_content = data.get('markdown')
        if not markdown_content:
            return {"error": "No 'markdown' field provided"}, 400

        file_type = data.get('fileType', 'pdf').lower()
        if file_type not in SUPPORTED_FILE_TYPES:
            return {
                "error": f"Invalid fileType '{file_type}'. Must be one of: {', '.join(sorted(SUPPORTED_FILE_TYPES))}"
            }, 400

        logger.info(
            "Received request to /generate with content-type: %s, fileType: %s, markdown length: %d characters",
            request.content_type,
            file_type,
            len(markdown_content),
        )
        
        # Create unique filenames
        run_id = str(uuid.uuid4())
        input_filename = f"slides_{run_id}.md"
        output_filename = f"slides_{run_id}.{file_type}"

        # Write Markdown file
        with open(input_filename, 'w') as f:
            f.write(markdown_content)

        # Run Marp
        cmd = [
            "marp",
            input_filename,
            MARP_FLAG[file_type],
            "--output", output_filename,
            "--allow-local-files",
        ]

        logger.info("Executing: %s", ' '.join(cmd))
        result = subprocess.run(cmd, check=True, capture_output=True, text=True)
        if result.stdout:
            logger.info("Marp stdout: %s", result.stdout.strip())
        if result.stderr:
            logger.warning("Marp stderr: %s", result.stderr.strip())

        return send_file(output_filename, mimetype=MIME_TYPE[file_type], as_attachment=True)

    except subprocess.CalledProcessError as e:
        logger.error("Marp Error: %s", e.stderr)
        return {"error": "Marp failed to generate file", "details": e.stderr}, 500

    except Exception as e:
        error_trace = traceback.format_exc()
        logger.exception("Server Error: %s", error_trace)
        return {"error": "Internal Server Error", "details": str(e), "trace": error_trace}, 500

    finally:
        if input_filename and os.path.exists(input_filename):
            os.remove(input_filename)


TEMPLATE_PATH = os.path.join(os.path.dirname(__file__), "template", "template.pptx")


@app.route("/v2/generate", methods=["POST"])
def generate_slides_v2():
    output_filename = None

    try:
        data = request.get_json()
        if not data:
            return {"error": "Invalid JSON body"}, 400

    except Exception as e:
        logger.error("Invalid JSON: %s", str(e))
        return {"error": "Invalid JSON body", "details": str(e)}, 400

    try:
        sections = data.get("sections")
        if not sections or not isinstance(sections, list):
            return {"error": "'sections' must be a non-empty array"}, 400

        validateSectionInput(sections) 

        logger.info(
            "POST /v2/generate – sections=%d, total_slides=%d",
            len(sections),
            sum(len(s.get("slides", [])) for s in sections),
        )

        run_id = str(uuid.uuid4())
        output_filename = f"slides_{run_id}.pptx"

        generate_pptx(TEMPLATE_PATH, sections, output_filename)

        with open(output_filename, "rb") as f:
            buf = io.BytesIO(f.read())
        buf.seek(0)

        return send_file(
            buf,
            mimetype="application/vnd.openxmlformats-officedocument.presentationml.presentation",
            as_attachment=True,
            download_name="presentation.pptx",
        )

    except Exception as e:
        error_trace = traceback.format_exc()
        logger.exception("Server Error: %s", error_trace)
        return {"error": "Internal Server Error"}, 500

    finally:
        if output_filename and os.path.exists(output_filename):
            os.remove(output_filename)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)