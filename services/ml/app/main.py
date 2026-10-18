"""FastAPI application entrypoint for the MonTra internal AI/ML service."""

from __future__ import annotations

from fastapi import FastAPI, Request
from fastapi.encoders import jsonable_encoder
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse

from app.api.routes import router as internal_router
from app.config import Settings
from app.schemas.prediction import ErrorDetail, ErrorResponse
from app.services.registry import ModelRegistry


def create_app(settings: Settings | None = None) -> FastAPI:
    active_settings = settings or Settings.from_env()
    registry = ModelRegistry(
        artifacts_dir=active_settings.artifacts_dir,
        category_model_version=active_settings.category_model_version,
    )

    app = FastAPI(
        title="MonTra Internal AI/ML Service",
        version="0.1.0",
        docs_url=None,
        redoc_url=None,
        openapi_url="/internal/v1/openapi.json",
    )
    app.state.settings = active_settings
    app.state.registry = registry

    @app.exception_handler(PermissionError)
    async def handle_permission_error(_request: Request, exc: PermissionError) -> JSONResponse:
        body = ErrorResponse(
            error=ErrorDetail(
                code="UNAUTHORIZED_INTERNAL_CALL",
                message=str(exc),
            )
        )
        return JSONResponse(status_code=401, content=body.model_dump(mode="json"))

    @app.exception_handler(RequestValidationError)
    async def handle_validation_error(
        _request: Request, exc: RequestValidationError
    ) -> JSONResponse:
        sanitized_errors = [
            {
                "loc": err.get("loc"),
                "msg": err.get("msg"),
                "type": err.get("type"),
            }
            for err in exc.errors()
        ]
        body = ErrorResponse(
            error=ErrorDetail(
                code="VALIDATION_ERROR",
                message="Payload permintaan ML tidak valid.",
                details={"errors": jsonable_encoder(sanitized_errors)},
            )
        )
        return JSONResponse(status_code=422, content=body.model_dump(mode="json"))

    app.include_router(internal_router)
    return app


app = create_app()
