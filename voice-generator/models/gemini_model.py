import io
import logging
import time
import wave

from google import genai
from google.genai import types

import config
from .base import TTSModel, TTSRequest
from .rate_limiter import SlidingWindowRateLimiter

logger = logging.getLogger(__name__)


class GeminiModel(TTSModel):
    """TTS backend using Google Gemini's text-to-speech API.

    Rate limiting is enforced locally via a sliding-window limiter so that
    we never exceed the Google-imposed request cap.  If Google still returns
    a 429 (e.g. due to shared quota), the call is retried with exponential
    backoff up to ``config.GEMINI_TTS_MAX_RETRIES`` times.
    """

    SAMPLE_RATE = 24_000

    _VOICES = [
        "Zephyr", "Puck", "Charon", "Kore", "Fenrir", "Leda",
        "Orus", "Aoede", "Callirrhoe", "Autonoe", "Enceladus", "Iapetus",
        "Umbriel", "Algieba", "Despina", "Erinome", "Algenib", "Rasalgethi",
        "Laomedeia", "Achernar", "Alnilam", "Schedar", "Gacrux", "Pulcherrima",
        "Achird", "Zubenelgenubi", "Vindemiatrix", "Sadachbia", "Sadaltager", "Sulafat",
    ]

    def __init__(self) -> None:
        self._client = genai.Client(api_key=config.GOOGLE_GENAI_API_KEY)
        self._rate_limiter = SlidingWindowRateLimiter(
            max_requests=config.GEMINI_TTS_MAX_REQUESTS,
            window_seconds=config.GEMINI_TTS_WINDOW_SECONDS,
        )

    # ------------------------------------------------------------------
    # TTSModel interface
    # ------------------------------------------------------------------

    def generate(self, request: TTSRequest) -> bytes:
        """Synthesise *request.text* via Gemini TTS and return WAV bytes."""
        voice = request.voice if request.voice != "default" else "Kore"

        # Block until we have a rate-limit slot.
        self._rate_limiter.acquire()

        logger.info(
            "Generating TTS with GeminiModel (voice=%s, language=%s, text length=%d)",
            voice, request.language, len(request.text),
        )
        pcm_data = self._call_with_retry(request.text, voice)
        return self._pcm_to_wav(pcm_data)

    def supported_voices(self) -> list[str]:
        return list(self._VOICES)

    def supported_languages(self) -> list[str]:
        return [
            "en", "th", "ja", "ko", "zh", "fr", "de", "es", "pt", "it",
            "nl", "pl", "ru", "tr", "ar", "hi", "id", "vi",
        ]

    # ------------------------------------------------------------------
    # Internal helpers
    # ------------------------------------------------------------------

    def _call_with_retry(self, text: str, voice: str) -> bytes:
        """Call Gemini TTS, retrying on 429 with exponential backoff."""
        last_exc: Exception | None = None
        max_retries = config.GEMINI_TTS_MAX_RETRIES
        base_delay = 2.0  # seconds

        for attempt in range(max_retries + 1):
            try:
                response = self._client.models.generate_content(
                    model=config.GEMINI_TTS_MODEL,
                    contents=text,
                    config=types.GenerateContentConfig(
                        response_modalities=["AUDIO"],
                        speech_config=types.SpeechConfig(
                            voice_config=types.VoiceConfig(
                                prebuilt_voice_config=types.PrebuiltVoiceConfig(
                                    voice_name=voice,
                                )
                            )
                        ),
                    ),
                )
                return response.candidates[0].content.parts[0].inline_data.data

            except Exception as exc:
                last_exc = exc
                if not self._is_rate_limit_error(exc):
                    raise
                if attempt < max_retries:
                    delay = base_delay * (2 ** attempt)
                    logger.warning(
                        "Gemini 429 rate-limited (attempt %d/%d), retrying in %.1fs",
                        attempt + 1, max_retries + 1, delay,
                    )
                    time.sleep(delay)

        raise RuntimeError(
            f"Gemini TTS failed after {max_retries + 1} attempts"
        ) from last_exc

    @staticmethod
    def _is_rate_limit_error(exc: Exception) -> bool:
        """Return True if *exc* represents an HTTP 429 from Google."""
        # google-genai wraps API errors; check status code and message.
        code = getattr(exc, "code", None) or getattr(exc, "status_code", None)
        if code == 429:
            return True
        if "429" in str(exc) or "RESOURCE_EXHAUSTED" in str(exc):
            return True
        return False

    def _pcm_to_wav(self, pcm_data: bytes) -> bytes:
        """Wrap raw PCM bytes from Gemini in a WAV header."""
        buf = io.BytesIO()
        with wave.open(buf, "wb") as wf:
            wf.setnchannels(1)
            wf.setsampwidth(2)
            wf.setframerate(self.SAMPLE_RATE)
            wf.writeframes(pcm_data)
        buf.seek(0)
        return buf.read()
