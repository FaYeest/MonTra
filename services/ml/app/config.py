"""Runtime configuration for the MonTra internal AI/ML service."""

from __future__ import annotations

import os
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class Settings:
    """Validated configuration loaded from environment variables."""

    app_env: str
    log_level: str
    host: str
    port: int
    artifacts_dir: Path
    category_model_version: str
    internal_service_token: str

    @classmethod
    def from_env(cls) -> "Settings":
        port_raw = os.getenv("ML_PORT", "8000").strip()
        try:
            port = int(port_raw)
        except ValueError as exc:
            raise ValueError("ML_PORT must be an integer") from exc
        if not (1 <= port <= 65535):
            raise ValueError("ML_PORT must be between 1 and 65535")

        artifacts_path = Path(os.getenv("ML_MODEL_ARTIFACTS_DIR", "artifacts")).resolve()
        return cls(
            app_env=os.getenv("APP_ENV", "development").strip() or "development",
            log_level=os.getenv("LOG_LEVEL", "INFO").strip().upper() or "INFO",
            host=os.getenv("ML_HOST", "0.0.0.0").strip() or "0.0.0.0",
            port=port,
            artifacts_dir=artifacts_path,
            category_model_version=os.getenv("ML_CATEGORY_MODEL_VERSION", "untrained-v0").strip()
            or "untrained-v0",
            internal_service_token=os.getenv("INTERNAL_SERVICE_TOKEN", "").strip(),
        )
