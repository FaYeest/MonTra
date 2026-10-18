"""Health and model readiness schemas matching packages/api-spec/openapi.yaml."""

from __future__ import annotations

from typing import Literal

from pydantic import BaseModel, Field


class ModelArtifactStatus(BaseModel):
    loaded: bool = Field(..., description="Whether an evaluated model artifact is loaded in memory")
    model_version: str = Field(..., description="Version identifier of the model")
    status: Literal["ready", "model_not_loaded"] = Field(
        ..., description="Current readiness state of the model artifact"
    )


class ActiveModelsMap(BaseModel):
    category_classifier: ModelArtifactStatus
    anomaly_detector: ModelArtifactStatus
    expense_forecaster: ModelArtifactStatus


class InternalMLHealthResponse(BaseModel):
    status: Literal["ok", "degraded"]
    service: str = "montra-ml"
    environment: str
    model_ready: bool
    active_models: ActiveModelsMap
    timestamp: str
