"""Application configuration.

All values are read from environment variables with sensible defaults so the
service works out-of-the-box locally while remaining fully configurable in
Docker / Kubernetes / cloud deployments.
"""
import os

# --------------------------------------------------------------------------
# Server
# --------------------------------------------------------------------------
IS_DEVELOPMENT: bool = os.getenv("ENV", "development").lower() == "development"
HOST: str = os.getenv("HOST", "0.0.0.0")
PORT: int = int(os.getenv("PORT", "8000"))
LOG_LEVEL: str = os.getenv("LOG_LEVEL", "info")

# --------------------------------------------------------------------------
# Concurrency
# --------------------------------------------------------------------------
# Maximum number of TTS synthesis jobs that may run simultaneously only for KokoroModel.
# Raise this only if the container has enough CPU/RAM for parallel inference.
MAX_CONCURRENT_TTS: int = int(os.getenv("MAX_CONCURRENT_TTS", "2"))

# --------------------------------------------------------------------------
# TTS defaults
# --------------------------------------------------------------------------
DEFAULT_MODEL: str = os.getenv("DEFAULT_TTS_MODEL", "kokoro")
DEFAULT_VOICE: str = os.getenv("DEFAULT_VOICE", "default")
DEFAULT_SPEED: float = float(os.getenv("DEFAULT_SPEED", "1.0"))
DEFAULT_LANGUAGE: str = os.getenv("DEFAULT_LANGUAGE", "a")

# --------------------------------------------------------------------------
# Local storage
# --------------------------------------------------------------------------
LOCAL_OUTPUT_DIR: str = os.getenv("LOCAL_OUTPUT_DIR", "/tmp/voice-output")

# --------------------------------------------------------------------------
# AWS / S3 storage
# --------------------------------------------------------------------------
S3_BUCKET: str = os.getenv("AWS_S3_BUCKET", "")
S3_PREFIX: str = os.getenv("AWS_S3_PREFIX", "voice-output/")
S3_REGION: str = os.getenv("AWS_S3_REGION", "us-east-1")
S3_PRESIGN_TTL: int = int(os.getenv("AWS_S3_PRESIGN_TTL", "3600"))

# Credentials — leave blank to use IAM role / instance profile
AWS_ACCESS_KEY_ID: str = os.getenv("AWS_ACCESS_KEY_ID", "")
AWS_SECRET_ACCESS_KEY: str = os.getenv("AWS_SECRET_ACCESS_KEY", "")

# --------------------------------------------------------------------------
# Google Gemini TTS
# --------------------------------------------------------------------------
GOOGLE_GENAI_API_KEY: str = os.getenv("GOOGLE_GENAI_API_KEY", "")
GEMINI_TTS_MODEL: str = os.getenv("GEMINI_TTS_MODEL", "gemini-2.5-flash-preview-tts")
# Sliding-window rate limit: max requests within the window (default: 10 / 120s).
GEMINI_TTS_MAX_REQUESTS: int = int(os.getenv("GEMINI_TTS_MAX_REQUESTS", "10"))
GEMINI_TTS_WINDOW_SECONDS: float = float(os.getenv("GEMINI_TTS_WINDOW_SECONDS", "120"))
# Retry attempts on 429 before raising an error.
GEMINI_TTS_MAX_RETRIES: int = int(os.getenv("GEMINI_TTS_MAX_RETRIES", "3"))
