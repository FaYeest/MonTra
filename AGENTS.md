# AGENTS.md — Engineering Guidelines for MonTra

This file defines non-negotiable architectural and engineering constraints for AI coding agents and human contributors working on **MonTra**.

## 1. Source of Truth Documents

- **Product Requirements:** [`docs/MonTra_PRD.md`](./docs/MonTra_PRD.md)
- **UI/UX Design Specification:** [`docs/design.md`](./docs/design.md)
- **Implementation Plan & Architecture:** [`docs/IMPLEMENTATION_PLAN.md`](./docs/IMPLEMENTATION_PLAN.md)
- **API Contracts:** [`packages/api-spec/openapi.yaml`](./packages/api-spec/openapi.yaml)

Always read the relevant sections of these documents before modifying domain logic, database schemas, API contracts, or Flutter screens.

## 2. Service Ownership & Boundaries

1. **Golang (`services/api`) is the authoritative backend.**
   - Owns authentication, resource-level authorization, accounts, categories, transactions, internal transfers, budgets, financial calculations, reports, and insight orchestration.
   - Is the **only** service permitted to write authoritative financial records to PostgreSQL.
2. **PostgreSQL is the single source of truth for financial records.**
   - All schema changes must go through numbered SQL migrations in `services/api/migrations/`.
   - Monetary values must use exact integer (`BIGINT` / `int64`) representation with explicit ISO 4217 currency codes (default `IDR`). Never use IEEE 754 floating-point types (`float32`/`float64`/`DOUBLE PRECISION`) for monetary storage or arithmetic.
   - Internal transfers between user accounts must execute atomically and must **never** be double-counted as external income or expense.
3. **Python FastAPI (`services/ml`) is an internal assistive AI/ML service.**
   - Exposes `/internal/v1/*` endpoints on the internal network only; never expose `services/ml` directly to public traffic.
   - Owns feature preprocessing, category classification inference, statistical anomaly detection, expense forecasting, and offline training/evaluation pipelines.
   - Must never fabricate ML predictions or confidence scores without an evaluated model artifact.
   - Outages, timeouts, or errors in `services/ml` must never block core transaction recording in `services/api` or `apps/mobile`.
4. **Flutter (`apps/mobile`) is the Android client.**
   - Follows `docs/design.md` tokens, typography, spacing, and 5-destination navigation (`Beranda`, `Transaksi`, `Tambah`, `Anggaran`, `Insight`, plus `Profil` in the top app bar).
   - Never computes authoritative balances or displays hardcoded sample numbers as real user financial data.

## 3. Security & Quality Rules

- **Strict Resource Authorization:** Scope every query on user data by the authenticated `user_id` from server context. Never trust a client-supplied `user_id`.
- **Secrets Hygiene:** Never commit `.env` files, real passwords, JWT secrets, or API keys. Keep `.env.example` updated with safe placeholders only.
- **Privacy & Log Redaction:** Never log passwords, tokens, or raw financial descriptions by default. Send only minimum required fields to `services/ml`.
- **Verification Before Claims:** Always execute formatters, linters, and automated tests (`go test ./...`, `pytest`/`unittest`, `flutter test`, `docker compose config`) and report actual results honestly.
