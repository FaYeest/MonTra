"""Internal API routes (`/internal/v1/*`) for the MonTra Python ML service."""

from __future__ import annotations

import hmac
from datetime import datetime, timezone

from fastapi import APIRouter, Depends, Header, Request
from fastapi.responses import JSONResponse

from app.config import Settings
from app.schemas.health import InternalMLHealthResponse
from app.schemas.prediction import (
    CategoryPredictionRequest,
    CategoryPredictionResponse,
    ErrorDetail,
    ErrorResponse,
)
from app.services.registry import ModelNotReadyError, ModelRegistry

router = APIRouter(prefix="/internal/v1", tags=["InternalML"])


def get_settings(request: Request) -> Settings:
    return request.app.state.settings


def get_registry(request: Request) -> ModelRegistry:
    return request.app.state.registry


def verify_internal_token(
    settings: Settings = Depends(get_settings),
    x_internal_service_token: str | None = Header(default=None),
) -> None:
    expected = settings.internal_service_token
    if not expected:
        return
    provided = x_internal_service_token or ""
    if not hmac.compare_digest(provided.encode("utf-8"), expected.encode("utf-8")):
        raise PermissionError("Invalid or missing X-Internal-Service-Token header.")


@router.get("/health", response_model=InternalMLHealthResponse)
def get_internal_health(
    settings: Settings = Depends(get_settings),
    registry: ModelRegistry = Depends(get_registry),
) -> InternalMLHealthResponse:
    models_status = registry.active_models_status()
    model_ready = registry.is_ready
    return InternalMLHealthResponse(
        status="ok" if model_ready else "degraded",
        service="montra-ml",
        environment=settings.app_env,
        model_ready=model_ready,
        active_models=models_status,
        timestamp=datetime.now(timezone.utc).isoformat(),
    )


@router.post(
    "/predict/category",
    response_model=CategoryPredictionResponse,
    responses={
        401: {"model": ErrorResponse},
        503: {"model": ErrorResponse},
    },
)
def predict_category(
    payload: CategoryPredictionRequest,
    _auth: None = Depends(verify_internal_token),
    registry: ModelRegistry = Depends(get_registry),
) -> CategoryPredictionResponse | JSONResponse:
    try:
        return registry.predict_category(payload)
    except ModelNotReadyError as exc:
        err = ErrorResponse(
            error=ErrorDetail(
                code="MODEL_NOT_READY",
                message=str(exc),
            )
        )
        return JSONResponse(status_code=503, content=err.model_dump())
