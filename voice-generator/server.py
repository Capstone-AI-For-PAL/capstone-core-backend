"""Voice-generator microservice.

Routes
------
GET  /health              — liveness probe
GET  /models              — list registered TTS models
POST /generate            — synthesise and stream WAV bytes directly
POST /generate/save       — synthesise, persist to a storage backend, return metadata
"""
import io
import uuid
import logging
from contextlib import asynccontextmanager
from typing import Optional

import uvicorn
from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse, JSONResponse
from pydantic import BaseModel, Field

import config
from models import get_model, _REGISTRY
from models.base import TTSRequest
from storage import get_storage
from storage.base import SaveResult

logging.basicConfig(level=config.LOG_LEVEL.upper())
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info(f"Is Enabled paid TTS: {config.ENABLE_PAID_TTS}")
    logger.info(f"Enabled TTS models: {list(_REGISTRY.keys())}")
    yield


app = FastAPI(
    title="Voice Generator Service",
    description="Text-to-speech microservice with pluggable models and storage backends.",
    version="1.0.0",
    lifespan=lifespan,
)


# ---------------------------------------------------------------------------
# Request / response schemas
# ---------------------------------------------------------------------------

class GenerateRequest(BaseModel):
    text: str = Field(..., description="Text to synthesise.")
    model: str = Field(config.DEFAULT_MODEL, description="TTS model name (e.g. 'kokoro').")
    voice: str = Field(config.DEFAULT_VOICE, description="Voice identifier (model-specific).")
    speed: float = Field(config.DEFAULT_SPEED, ge=0.1, le=4.0, description="Speech rate multiplier.")
    language: str = Field(
        config.DEFAULT_LANGUAGE,
        description="Language/dialect code (model-specific, e.g. 'a' = American English for Kokoro).",
    )


class SaveRequest(GenerateRequest):
    storage: str = Field("local", description="Storage backend: 'local' or 's3'.")
    filename: Optional[str] = Field(
        None,
        description="Output filename. Auto-generated UUID if omitted.",
    )


class SaveResponse(BaseModel):
    location: str
    url: Optional[str]
    content_type: str
    backend: str
    filename: str


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _build_tts_request(req: GenerateRequest) -> TTSRequest:
    return TTSRequest(
        text=req.text,
        voice=req.voice,
        speed=req.speed,
        language=req.language,
    )


def _synthesise(req: GenerateRequest) -> bytes:
    """Run TTS synthesis and return raw WAV bytes.

    Concurrency / rate-limit controls are managed inside each model
    implementation (e.g. semaphore in KokoroModel, sliding-window in
    GeminiModel), so this helper simply delegates.
    """
    try:
        model = get_model(req.model)
    except ValueError as exc:
        raise HTTPException(status_code=400, detail=str(exc))

    try:
        return model.generate(_build_tts_request(req))
    except Exception as exc:
        logger.exception("TTS generation failed")
        raise HTTPException(status_code=500, detail=f"TTS generation failed: {exc}")


# ---------------------------------------------------------------------------
# Routes
# ---------------------------------------------------------------------------

@app.get("/health", tags=["ops"])
def health() -> dict:
    """Liveness probe."""
    return {"status": "ok"}


@app.get("/models", tags=["ops"])
def list_models() -> dict:
    """Return registered TTS model names."""
    return {"models": list(_REGISTRY.keys())}


@app.post(
    "/generate",
    tags=["voice"],
    response_class=StreamingResponse,
    responses={
        200: {
            "content": {"audio/wav": {}},
            "description": "WAV audio file streamed directly.",
        }
    },
)
def generate(req: GenerateRequest):
    """Synthesise text and stream the resulting WAV file back to the caller.

    The response body is the raw WAV file; set ``Accept: audio/wav`` or simply
    read the binary content.
    """
    wav_bytes = _synthesise(req)

    filename = f"{uuid.uuid4()}.wav"
    headers = {"Content-Disposition": f'attachment; filename="{filename}"'}

    return StreamingResponse(
        io.BytesIO(wav_bytes),
        media_type="audio/wav",
        headers=headers,
    )


@app.post("/generate/save", response_model=SaveResponse, tags=["voice"])
def generate_and_save(req: SaveRequest):
    """Synthesise text, persist the WAV to the chosen storage backend, and
    return a JSON payload describing where the file was stored.

    - ``storage="local"`` saves to the local filesystem (``LOCAL_OUTPUT_DIR``).
    - ``storage="s3"``   uploads to the configured S3 bucket and returns a
      pre-signed URL valid for ``S3_PRESIGN_TTL`` seconds.
    """
    wav_bytes = _synthesise(req)

    filename = req.filename or f"{uuid.uuid4()}.wav"
    if not filename.endswith(".wav"):
        filename += ".wav"

    try:
        backend = get_storage(req.storage)
    except ValueError as exc:
        raise HTTPException(status_code=400, detail=str(exc))

    try:
        result: SaveResult = backend.save(wav_bytes, filename)
    except Exception as exc:
        logger.exception("Storage save failed")
        raise HTTPException(status_code=500, detail=f"Storage save failed: {exc}")

    return SaveResponse(
        location=result.location,
        url=result.url,
        content_type=result.content_type,
        backend=result.backend,
        filename=filename,
    )


# ---------------------------------------------------------------------------
# Entrypoint
# ---------------------------------------------------------------------------

if __name__ == "__main__":
    uvicorn.run(
        "server:app",
        host=config.HOST,
        port=config.PORT,
        log_level=config.LOG_LEVEL,
        reload=config.IS_DEVELOPMENT,
    )
