from .base import TTSModel, TTSRequest
from .kokoro_model import KokoroModel

_REGISTRY: dict[str, type[TTSModel]] = {
    "kokoro": KokoroModel,
}

# Singleton cache — one instance per model name, reused across all requests.
# This ensures the pipeline and HuggingFace weight downloads only happen once.
_INSTANCES: dict[str, TTSModel] = {}


def get_model(name: str) -> TTSModel:
    """Return a shared TTSModel instance by registry name.

    Instances are created once and reused for the lifetime of the process,
    avoiding redundant HuggingFace Hub checks and pipeline re-initialisations
    on every request.
    """
    if name not in _INSTANCES:
        cls = _REGISTRY.get(name)
        if cls is None:
            raise ValueError(
                f"Unknown TTS model '{name}'. Available: {list(_REGISTRY.keys())}"
            )
        _INSTANCES[name] = cls()
    return _INSTANCES[name]
