"""Thread-safe sliding-window rate limiter.

Used to enforce external API rate limits (e.g. Google Gemini TTS) without
dropping requests.  Callers that exceed the window simply block until a
slot becomes available.
"""
import threading
import time
from collections import deque


class SlidingWindowRateLimiter:
    """Allow at most *max_requests* calls within a rolling *window_seconds* window.

    When the limit is reached, :meth:`acquire` blocks the calling thread until
    the oldest request ages out of the window.  This gives natural FIFO
    queuing behaviour without returning errors.

    Parameters
    ----------
    max_requests:
        Maximum number of requests permitted inside the window.
    window_seconds:
        Length of the sliding window in seconds.
    """

    def __init__(self, max_requests: int, window_seconds: float) -> None:
        self._max_requests = max_requests
        self._window_seconds = window_seconds
        self._timestamps: deque[float] = deque()
        self._lock = threading.Lock()

    def acquire(self) -> None:
        """Block until a rate-limit slot is available, then record the call."""
        while True:
            with self._lock:
                now = time.monotonic()
                # Evict timestamps that have fallen outside the window.
                while self._timestamps and now - self._timestamps[0] >= self._window_seconds:
                    self._timestamps.popleft()

                if len(self._timestamps) < self._max_requests:
                    self._timestamps.append(now)
                    return  # Slot acquired.

                # Calculate how long to sleep until the oldest entry expires.
                sleep_for = self._window_seconds - (now - self._timestamps[0])

            # Sleep *outside* the lock so other threads aren't blocked.
            time.sleep(max(sleep_for, 0.01))
