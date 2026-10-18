"""Request, response, and error schemas for internal ML prediction endpoints."""

from __future__ import annotations

from typing import Any, Literal

from pydantic import BaseModel, Field, field_validator


class CategoryPredictionRequest(BaseModel):
    description: str = Field(..., min_length=1, max_length=255)
    transaction_type: Literal["income", "expense"]
    currency: str = Field(..., min_length=3, max_length=3)
    amount: int = Field(..., gt=0, description="Positive integer amount in whole currency units")

    @field_validator("description")
    @classmethod
    def validate_non_blank_description(cls, value: str) -> str:
        cleaned = value.strip()
        if not cleaned:
            raise ValueError("description must not be blank")
        return cleaned

    @field_validator("currency")
    @classmethod
    def validate_currency_code(cls, value: str) -> str:
        cleaned = value.strip().upper()
        if len(cleaned) != 3 or not cleaned.isalpha():
            raise ValueError("currency must be a 3-letter ISO 4217 code")
        return cleaned


class CategoryPredictionResponse(BaseModel):
    predicted_category: str
    confidence: float = Field(..., ge=0.0, le=1.0)
    model_version: str


class ErrorDetail(BaseModel):
    code: str
    message: str
    details: dict[str, Any] | None = None


class ErrorResponse(BaseModel):
    error: ErrorDetail
