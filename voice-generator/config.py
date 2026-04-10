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
# Maximum number of TTS synthesis jobs that may run simultaneously.
# Raise this only if the container has enough CPU/RAM for parallel inference.
MAX_CONCURRENT_TTS: int = int(os.getenv("MAX_CONCURRENT_TTS", "2"))

# --------------------------------------------------------------------------
# TTS defaults
# --------------------------------------------------------------------------
DEFAULT_MODEL: str = os.getenv("DEFAULT_TTS_MODEL", "kokoro")
DEFAULT_VOICE: str = os.getenv("DEFAULT_VOICE", "af_heart")
DEFAULT_SPEED: float = float(os.getenv("DEFAULT_SPEED", "1.0"))
DEFAULT_LANGUAGE: str = os.getenv("DEFAULT_LANGUAGE", "a")

# --------------------------------------------------------------------------
# Local storage
# --------------------------------------------------------------------------
LOCAL_OUTPUT_DIR: str = os.getenv("LOCAL_OUTPUT_DIR", "/tmp/voice-output")

# --------------------------------------------------------------------------
# AWS / S3 storage
# --------------------------------------------------------------------------
S3_BUCKET: str = os.getenv("S3_BUCKET", "")
S3_PREFIX: str = os.getenv("S3_PREFIX", "voice-output/")
S3_REGION: str = os.getenv("S3_REGION", "us-east-1")
S3_PRESIGN_TTL: int = int(os.getenv("S3_PRESIGN_TTL", "3600"))

# Credentials — leave blank to use IAM role / instance profile
AWS_ACCESS_KEY_ID: str = os.getenv("AWS_ACCESS_KEY_ID", "")
AWS_SECRET_ACCESS_KEY: str = os.getenv("AWS_SECRET_ACCESS_KEY", "")
