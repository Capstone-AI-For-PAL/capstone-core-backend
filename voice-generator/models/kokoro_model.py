import io
import threading
import numpy as np
import soundfile as sf
from kokoro import KPipeline

import config
from .base import TTSModel, TTSRequest


class KokoroModel(TTSModel):
    """TTS backend powered by the Kokoro neural TTS library.

    Model weights (~300 MB) are downloaded automatically on first use and
    cached by the kokoro library.

    Concurrency is capped by a semaphore (``config.MAX_CONCURRENT_TTS``)
    because each inference run is CPU/memory-heavy.  Extra callers block
    until a slot is free.
    """

    # Kokoro's native output sample rate
    SAMPLE_RATE = 24_000

    # Voices bundled with Kokoro at the time of writing
    _VOICES = [
        "af_heart", "af_bella", "af_sarah", "af_sky",
        "am_adam", "am_michael",
        "bf_emma", "bf_isabella",
        "bm_george", "bm_lewis",
    ]

    def __init__(self) -> None:
        # Pipeline is initialised lazily so the import doesn't block startup
        self._pipeline: KPipeline | None = None
        # Serialise all pipeline access: KPipeline is not thread-safe and
        # concurrent inference causes crashes / resource exhaustion.
        self._lock = threading.Lock()
        self._semaphore = threading.Semaphore(config.MAX_CONCURRENT_TTS)

    def _get_pipeline(self, lang_code: str) -> KPipeline:
        """Return (and lazily initialise) the Kokoro pipeline.

        Must be called while holding ``self._lock``.
        """
        if self._pipeline is None or self._pipeline.lang_code != lang_code:
            self._pipeline = KPipeline(lang_code=lang_code)
        return self._pipeline

    # ------------------------------------------------------------------
    # TTSModel interface
    # ------------------------------------------------------------------

    def generate(self, request: TTSRequest) -> bytes:
        """Synthesise *request.text* and return a WAV file as bytes."""
        self._semaphore.acquire()
        try:
            with self._lock:
                pipeline = self._get_pipeline(request.language)

                voice = request.voice if request.voice != "default" else "af_heart"

                generator = pipeline(
                    request.text,
                    voice=voice,
                    speed=request.speed,
                    split_pattern=r"\n+",
                )

                all_audio: list[np.ndarray] = []
                for _graphemes, _phonemes, audio_chunk in generator:
                    all_audio.append(audio_chunk)

            if not all_audio:
                raise RuntimeError("Kokoro returned no audio chunks.")

            combined = np.concatenate(all_audio)

            # Encode to WAV in-memory
            buf = io.BytesIO()
            sf.write(buf, combined, self.SAMPLE_RATE, format="WAV")
            buf.seek(0)
            return buf.read()
        finally:
            self._semaphore.release()

    def supported_voices(self) -> list[str]:
        return list(self._VOICES)

    def supported_languages(self) -> list[str]:
        # 'a' = American English, 'b' = British English
        # Future Kokoro releases may add more.
        return ["a", "b"]
