# Model Evaluation Harness (`services/ml/evaluation`)

This directory houses reproducible model evaluation scripts and metric reports (scheduled for Phase 3 & Phase 4 per `docs/IMPLEMENTATION_PLAN.md`).

## Required Evaluation Gates

1. **Transaction Categorization (`AI-01`):**
   - Primary metric: **Macro-F1** on held-out test set (development target >= 0.80).
   - Required diagnostics: per-category precision/recall, confusion matrix, rare-category performance, and confidence calibration check.
2. **Spending Anomaly Detection (`AI-02`):**
   - Evaluate precision per alert and false-positive rate against historical baselines.
3. **Monthly Expense Forecasting (`AI-03`):**
   - Evaluate **MAE** and **WAPE** against actual end-of-month totals and compare directly against a simple historical average baseline.
   - Enforce cold-start thresholds when historical coverage is insufficient.
4. **Honest Reporting:** Never claim model accuracy or confidence calibration without an executed evaluation report generated from held-out data.
