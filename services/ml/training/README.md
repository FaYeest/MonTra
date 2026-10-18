# Model Training Pipeline (`services/ml/training`)

This directory houses the offline training scripts for MonTra's assistive ML models (scheduled for Phase 3 & Phase 4 per `docs/IMPLEMENTATION_PLAN.md`).

## Architectural & Data Rules

1. **Separation from Inference API:** Training jobs never run inside the request path of the FastAPI inference container (`app/`).
2. **Leakage Prevention:** Datasets must be split into deterministic `train`, `validation`, and `test` partitions before fitting any vectorizer (`TfidfVectorizer`), scaler, or classifier.
3. **Baseline First (`AI-01`):**
   - Start with `TF-IDF + LogisticRegression` from `scikit-learn`.
   - Normalize merchant/description text via `app.models.metadata.normalize_merchant_description`.
   - Do not rely on transaction `amount` as a sole predictor of category.
4. **Artifact Serialization:** Save trained pipelines using `joblib` with explicit version metadata (`category-v1.joblib`) only after passing the evaluation gates in `services/ml/evaluation/`.
