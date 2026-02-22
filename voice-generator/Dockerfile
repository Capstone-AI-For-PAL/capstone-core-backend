# ---------------------------------------------------------------------------
# Stage 1 — dependency builder
# ---------------------------------------------------------------------------
FROM python:3.11-slim AS builder

WORKDIR /build

# System libs needed to compile Python packages (soundfile → libsndfile)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libsndfile1-dev \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --upgrade pip \
 && pip install --prefix=/install --no-cache-dir -r requirements.txt \
 && PYTHONPATH=/install/lib/python3.11/site-packages \
    pip install --prefix=/install --no-cache-dir \
    https://github.com/explosion/spacy-models/releases/download/en_core_web_sm-3.8.0/en_core_web_sm-3.8.0-py3-none-any.whl


# ---------------------------------------------------------------------------
# Stage 2 — runtime image
# ---------------------------------------------------------------------------
FROM python:3.11-slim AS runtime

# espeak-ng: required by the Kokoro phonemiser
# libsndfile1: runtime shared library used by soundfile
RUN apt-get update && apt-get install -y --no-install-recommends \
    espeak-ng \
    libsndfile1 \
    && rm -rf /var/lib/apt/lists/*

# Copy installed Python packages from builder (includes en_core_web_sm)
COPY --from=builder /install /usr/local

# ---------------------------------------------------------------------------
# Application setup
# ---------------------------------------------------------------------------
WORKDIR /app

COPY config.py     ./config.py
COPY server.py     ./server.py
COPY models/       ./models/
COPY storage/      ./storage/

# Kokoro downloads model weights on first use; point the cache to a
# predictable path so a volume or bind mount can persist the weights
# between container restarts and avoid repeated downloads.
ENV HF_HOME=/app/.cache/huggingface
ENV KOKORO_CACHE=/app/.cache/kokoro
RUN mkdir -p /app/.cache

# Create the default local output directory
RUN mkdir -p /tmp/voice-output

# ---------------------------------------------------------------------------
# Runtime configuration (all overridable via docker run -e or env file)
# ---------------------------------------------------------------------------
ENV HOST=0.0.0.0
ENV PORT=8000
ENV LOG_LEVEL=info
ENV DEFAULT_TTS_MODEL=kokoro
ENV DEFAULT_VOICE=af_heart
ENV DEFAULT_SPEED=1.0
ENV DEFAULT_LANGUAGE=a
ENV LOCAL_OUTPUT_DIR=/tmp/voice-output

EXPOSE 8000

# Run as non-root
RUN useradd -m appuser && chown -R appuser /app /tmp/voice-output
USER appuser

CMD ["python", "-m", "uvicorn", "server:app", "--host", "0.0.0.0", "--port", "8000"]
