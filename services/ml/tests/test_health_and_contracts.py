"""Unit and contract tests for the MonTra internal AI/ML service."""

from __future__ import annotations

import unittest
from pathlib import Path

from fastapi.testclient import TestClient

from app.config import Settings
from app.main import create_app
from app.models.metadata import normalize_merchant_description
from app.schemas.prediction import CategoryPredictionRequest


class TestMLServiceFoundation(unittest.TestCase):
    def setUp(self) -> None:
        self.settings = Settings(
            app_env="test",
            log_level="INFO",
            host="127.0.0.1",
            port=8000,
            artifacts_dir=Path("artifacts"),
            category_model_version="untrained-v0",
            internal_service_token="test-internal-secret",
        )
        self.app = create_app(self.settings)
        self.client = TestClient(self.app)

    def test_normalize_merchant_description(self) -> None:
        raw = "  KOPI KENANGAN #042 - Jakarta!! "
        self.assertEqual(
            normalize_merchant_description(raw),
            "kopi kenangan 042 jakarta",
        )

    def test_health_endpoint_reports_honest_untrained_state(self) -> None:
        response = self.client.get("/internal/v1/health")
        self.assertEqual(response.status_code, 200)
        body = response.json()
        self.assertEqual(body["service"], "montra-ml")
        self.assertEqual(body["status"], "degraded")
        self.assertFalse(body["model_ready"])
        self.assertFalse(body["active_models"]["category_classifier"]["loaded"])
        self.assertEqual(
            body["active_models"]["category_classifier"]["status"],
            "model_not_loaded",
        )

    def test_predict_category_requires_internal_service_token(self) -> None:
        payload = {
            "description": "Kopi Kenangan",
            "transaction_type": "expense",
            "currency": "IDR",
            "amount": 25000,
        }
        response = self.client.post("/internal/v1/predict/category", json=payload)
        self.assertEqual(response.status_code, 401)
        self.assertEqual(response.json()["error"]["code"], "UNAUTHORIZED_INTERNAL_CALL")

    def test_predict_category_refuses_fake_predictions_when_untrained(self) -> None:
        payload = {
            "description": "Kopi Kenangan",
            "transaction_type": "expense",
            "currency": "IDR",
            "amount": 25000,
        }
        response = self.client.post(
            "/internal/v1/predict/category",
            json=payload,
            headers={"X-Internal-Service-Token": "test-internal-secret"},
        )
        self.assertEqual(response.status_code, 503)
        body = response.json()
        self.assertEqual(body["error"]["code"], "MODEL_NOT_READY")

    def test_predict_category_validates_input_contract(self) -> None:
        invalid_payload = {
            "description": "   ",
            "transaction_type": "expense",
            "currency": "IDR",
            "amount": -100,
        }
        response = self.client.post(
            "/internal/v1/predict/category",
            json=invalid_payload,
            headers={"X-Internal-Service-Token": "test-internal-secret"},
        )
        self.assertEqual(response.status_code, 422)
        self.assertEqual(response.json()["error"]["code"], "VALIDATION_ERROR")

    def test_schema_normalizes_currency(self) -> None:
        req = CategoryPredictionRequest(
            description="  Warung Makan  ",
            transaction_type="expense",
            currency="idr",
            amount=18000,
        )
        self.assertEqual(req.description, "Warung Makan")
        self.assertEqual(req.currency, "IDR")


if __name__ == "__main__":
    unittest.main()
