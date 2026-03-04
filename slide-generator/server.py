from flask import Flask, request, send_file
import subprocess
import os
import uuid
import traceback

app = Flask(__name__)

SUPPORTED_FILE_TYPES = {"pdf", "pptx"}
MARP_FLAG = {"pdf": "--pdf", "pptx": "--pptx"}
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

        print(f"Received request to /generate with content-type: {request.content_type}, fileType: {file_type}, markdown length: {len(markdown_content)} characters")
        
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

        print(f"Executing: {' '.join(cmd)}")
        subprocess.run(cmd, check=True, capture_output=True, text=True)

        return send_file(output_filename, mimetype=MIME_TYPE[file_type], as_attachment=True)

    except subprocess.CalledProcessError as e:
        print(f"Marp Error: {e.stderr}")
        return {"error": "Marp failed to generate file", "details": e.stderr}, 500

    except Exception as e:
        error_trace = traceback.format_exc()
        print(f"Server Error: {error_trace}")
        return {"error": "Internal Server Error", "details": str(e), "trace": error_trace}, 500

    finally:
        if input_filename and os.path.exists(input_filename):
            os.remove(input_filename)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)