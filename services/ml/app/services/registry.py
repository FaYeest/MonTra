"""Model registry and honest readiness enforcement.

During Phase 0 (before training and evaluation in Phase 3/4), no trained model artifact is bundled.
The registry explicitly reports `loaded=False` and refuses to fabricate predictions or confidence scores.
"""

from __future__ import annotations

from pathlib import Path

from app.models.metadata import ModelMetadata, normalize_merchant_description
from app.schemas.health import ActiveModelsMap, ModelArtifactStatus
from app.schemas.prediction import CategoryPredictionRequest, CategoryPredictionResponse


class ModelNotReadyError(RuntimeError):
    """Raised when an inference request arrives before an evaluated model artifact is loaded."""


class ModelRegistry:
    """Manages loaded model artifacts and exposes honest readiness states."""

    def __init__(self, artifacts_dir: Path, category_model_version: str) -> None:
        self._artifacts_dir = artifacts_dir
        self._category_model_version = category_model_version
        self._category_artifact_path = artifacts_dir / f"{category_model_version}.joblib"

    @property
    def category_metadata(self) -> ModelMetadata:
        is_loaded = (
            self._category_model_version != "untrained-v0"
            and self._category_artifact_path.is_file()
        )
        return ModelMetadata(
            name="category_classifier",
            model_version=self._category_model_version,
            loaded=is_loaded,
        )

    def active_models_status(self) -> ActiveModelsMap:
        cat_meta = self.category_metadata
        return ActiveModelsMap(
            category_classifier=ModelArtifactStatus(
                loaded=cat_meta.loaded,
                model_version=cat_meta.model_version,
                status="ready" if cat_meta.loaded else "model_not_loaded",
            ),
            anomaly_detector=ModelArtifactStatus(
                loaded=False,
                model_version="untrained-v0",
                status="model_not_loaded",
            ),
            expense_forecaster=ModelArtifactStatus(
                loaded=False,
                model_version="untrained-v0",
                status="model_not_loaded",
            ),
        )

    @property
    def is_ready(self) -> bool:
        return self.category_metadata.loaded

    def predict_category(
        self, request: CategoryPredictionRequest
    ) -> CategoryPredictionResponse:
        normalized_desc = normalize_merchant_description(request.description)
        if not normalized_desc:
            raise ValueError("description contains no normalizable alphanumeric tokens")

        if not self.category_metadata.loaded:
            raise ModelNotReadyError(
                "Model kategorisasi belum dimuat atau belum melalui evaluasi."
            )

        # Phase 3 will load the evaluated scikit-learn pipeline from self._category_artifact_path.
        raise ModelNotReadyError(
            "Model kategorisasi belum dimuat atau belum melalui evaluasi."
        )
