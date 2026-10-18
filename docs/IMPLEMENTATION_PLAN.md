# MonTra — Architecture and Implementation Plan

**Status:** Baseline Architecture & Phase 0 Foundation Plan  
**Source Requirements:** [`MonTra_PRD.md`](./MonTra_PRD.md) (v1.1)  
**UI/UX Specification:** [`design.md`](./design.md) (v1.0)  
**Target Platform (MVP):** Android (Flutter / Dart), Bahasa Indonesia, IDR (`Rp`) default currency

---

## 1. Product Architecture Summary

MonTra (Money Tracker) is a personal finance tracking application with assistive, explainable AI/ML capabilities tailored initially for Indonesian users. The system is designed around a strict separation between **deterministic financial record-keeping** and **probabilistic/statistical AI/ML assistance**:

```text
+---------------------------------------------------------------+
|                Flutter Mobile App (Android)                   |
|  UI/UX tokens (design.md), IDR formatting, state management,  |
|  form validation, resilient offline/empty/error/ML states     |
+---------------------------------------------------------------+
                                |
                                | HTTPS / JSON REST (/api/v1/*)
                                v
+---------------------------------------------------------------+
|               Reverse Proxy (Caddy / Nginx)                   |
|        Public entrypoint, TLS termination, rate limits        |
+---------------------------------------------------------------+
                                |
                                v
+---------------------------------------------------------------+
|                Golang Primary Backend (services/api)          |
|  Authoritative owner of auth, users, accounts, transactions,  |
|  internal transfers, categories, budgets, reports, insights,  |
|  and resilient orchestration of internal ML requests          |
+---------------------------------------------------------------+
             |                                      |
             | SQL (pgx, atomic DB tx)              | Internal HTTP/JSON
             v                                      | (/internal/v1/*)
+-------------------------------+                   v
|      PostgreSQL Database      |   +-----------------------------------+
|  Authoritative source of      |   | Python AI/ML Service (services/ml)|
|  truth for all financial      |   | FastAPI, Pydantic, scikit-learn,  |
|  records, balances, budgets,  |   | pandas, NumPy. Category inference,|
|  audit metadata & predictions |   | anomaly detection, forecasting    |
+-------------------------------+   +-----------------------------------+
```

### Core Architectural Invariants

1. **Financial Correctness First:** All account balances, cash-flow summaries, budget utilization percentages, and historical comparisons are computed deterministically by the Golang backend and PostgreSQL. Neither the mobile presentation layer nor any AI/ML model is permitted to compute or mutate authoritative balances.
2. **Authoritative Single Writer:** PostgreSQL is the sole source of truth for financial data, accessed exclusively by the Golang backend.
3. **Internal-Only AI/ML Service:** The Python FastAPI service operates strictly on an internal network (`/internal/v1/*`). It is never exposed to the public internet and has no direct write access to authoritative financial tables.
4. **Graceful Degradation:** Every core user workflow (recording income, expenses, or transfers; managing accounts, categories, and budgets; viewing reports) remains 100% operational when the Python AI/ML service is slow, disabled, or unreachable.
5. **Explainability & User Control:** Every AI category suggestion can be overridden in one tap. Every statistical insight or forecast must cite its time period, numeric evidence, comparison baseline, and uncertainty/limitations.

---

## 2. Service Responsibilities and Boundaries

| Component | Technology | Owns / Responsibilities | Explicit Prohibitions |
|---|---|---|---|
| **Mobile App** (`apps/mobile`) | Flutter, Dart | Design system tokens (`design.md`), 5-destination navigation (`Beranda`, `Transaksi`, `Tambah`, `Anggaran`, `Insight`) + top-bar `Profil`, input validation, IDR locale formatting, accessibility semantics (48×48 dp targets, TalkBack labels), loading/empty/error/ML-unavailable states. | Must not compute authoritative account balances, silently apply low-confidence ML categories, display fake data as real records, or use decorative neon/AI gradients. |
| **Primary Backend** (`services/api`) | Go (`net/http` / Chi, `pgx`, `log/slog`) | Public `/api/v1` REST API, authentication, resource-level authorization (`user_id` scoping), atomic database transactions, internal transfers, deterministic balance & budget calculations, CSV export, insight evidence verification, timeout-bounded `mlclient`. | Must not block transaction creation on ML failures, trust client-supplied `user_id`, expose stack traces, or use floating-point types (`float32`/`float64`) for monetary arithmetic. |
| **Database** (`postgres`) | PostgreSQL 16 | Relational integrity, foreign keys, check constraints (non-negative amounts, distinct transfer accounts), transactional consistency, indexed queries per user/date/account/category, schema versioning via SQL migrations. | Must not store plaintext passwords or expose ports publicly in production deployments. |
| **AI/ML Service** (`services/ml`) | Python, FastAPI, Pydantic, scikit-learn, pandas, NumPy, Joblib | Feature normalization, transaction category classification (`/internal/v1/predict/category`), statistical spending anomaly detection (`/internal/v1/detect/anomalies`), monthly expense forecasting (`/internal/v1/predict/expense-forecast`), offline training & evaluation pipelines. | Must not be publicly routable, directly mutate PostgreSQL financial records, fabricate confidence scores without evaluated models, or label anomalies as "fraud". |
| **API Specification** (`packages/api-spec`) | OpenAPI 3.1 (`openapi.yaml`) | Contract definitions for public `/api/v1/*` endpoints and internal `/internal/v1/*` endpoints, shared error envelopes, pagination schemas, and validation constraints. | Must not diverge from implemented request/response schemas across Go, Python, and Flutter. |
| **Reverse Proxy** (`infrastructure/proxy`) | Caddy / Nginx | Routing external traffic exclusively to `services/api`, TLS termination in non-local environments, security headers. | Must never proxy external requests to `services/ml`. |

---

## 3. Repository Structure

Adapted to preserve existing documentation in `docs/MonTra_PRD.md` and `docs/design.md`:

```text
montra/
├── AGENTS.md
├── README.md
├── .env.example
├── .gitignore
├── docker-compose.yml
├── docs/
│   ├── MonTra_PRD.md              # Existing PRD (source of truth for product scope)
│   ├── design.md                  # Existing UI/UX specification
│   └── IMPLEMENTATION_PLAN.md     # Architecture & phased implementation plan
├── apps/
│   └── mobile/                    # Flutter Android application
│       ├── analysis_options.yaml
│       ├── pubspec.yaml
│       ├── lib/
│       │   ├── main.dart
│       │   ├── app/
│       │   │   ├── app.dart
│       │   │   ├── router.dart
│       │   │   └── theme/
│       │   │       ├── app_colors.dart
│       │   │       ├── app_shapes.dart
│       │   │       ├── app_spacing.dart
│       │   │       ├── app_theme.dart
│       │   │       └── app_typography.dart
│       │   ├── core/
│       │   │   ├── formatting/
│       │   │   │   └── idr_formatter.dart
│       │   │   ├── networking/
│       │   │   │   └── api_config.dart
│       │   │   └── widgets/
│       │   │       └── montra_app_shell.dart
│       │   └── features/
│       │       ├── home/
│       │       ├── transactions/
│       │       ├── budgets/
│       │       ├── insights/
│       │       └── profile/
│       └── test/
│           ├── idr_formatter_test.dart
│           └── app_shell_test.dart
├── services/
│   ├── api/                       # Golang authoritative backend
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── cmd/
│   │   │   └── api/
│   │   │       └── main.go
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── database/
│   │   │   ├── domain/
│   │   │   │   └── money/
│   │   │   ├── httpserver/
│   │   │   ├── mlclient/
│   │   │   ├── auth/              # Phase 1 module boundary
│   │   │   ├── users/             # Phase 1 module boundary
│   │   │   ├── accounts/          # Phase 1 module boundary
│   │   │   ├── categories/        # Phase 1 module boundary
│   │   │   ├── transactions/      # Phase 1 module boundary
│   │   │   ├── transfers/         # Phase 1 module boundary
│   │   │   ├── budgets/           # Phase 2 module boundary
│   │   │   ├── reports/           # Phase 2 module boundary
│   │   │   ├── insights/          # Phase 4 module boundary
│   │   │   └── notifications/     # Phase 2/4 module boundary
│   │   └── migrations/
│   │       ├── 000001_init_foundation.up.sql
│   │       ├── 000001_init_foundation.down.sql
│   │       └── README.md
│   └── ml/                        # Python FastAPI internal AI/ML service
│       ├── Dockerfile
│       ├── pyproject.toml
│       ├── app/
│       │   ├── __init__.py
│       │   ├── main.py
│       │   ├── config.py
│       │   ├── api/
│       │   │   └── routes.py
│       │   ├── schemas/
│       │   │   ├── health.py
│       │   │   └── prediction.py
│       │   ├── services/
│       │   │   └── registry.py
│       │   └── models/
│       │       └── metadata.py
│       ├── training/
│       │   └── README.md
│       ├── evaluation/
│       │   └── README.md
│       ├── artifacts/
│       │   └── .gitkeep
│       └── tests/
│           └── test_health_and_contracts.py
├── packages/
│   └── api-spec/
│       └── openapi.yaml           # Unified OpenAPI 3.1 contract
└── infrastructure/
    ├── docker/
    │   └── postgres-init.sql
    └── proxy/
        └── Caddyfile
```

---

## 4. Major Implementation Phases and Dependencies

### Phase 0 — Architecture & Foundation Setup (Current Iteration)
- **Scope:** Repository structure, implementation plan, OpenAPI baseline specification, Docker Compose environment, Go API skeleton with configuration/database health/ML client/money primitive, Python FastAPI skeleton with health check and model readiness contracts, Flutter foundation with `design.md` tokens, IDR formatter, and 5-tab navigation shell.
- **Dependencies:** None.

### Phase 1 — Core Financial Ledger & Authentication
- **Scope:**
  - PostgreSQL migrations for `users`, `accounts`, `categories` (system seeded + user custom), `transactions`, and `transfers`.
  - Authentication (`register`, `login`, `logout`, password reset, session/JWT middleware) and `me` profile endpoints (`FR-01`).
  - Virtual account management and opening balance tracking (`FR-05`).
  - Category management (`FR-06`).
  - Transaction CRUD (`income`, `expense`) with idempotency keys, pagination, and filters (`FR-03`).
  - Atomic internal transfers (`source_account_id` -> `destination_account_id` + optional `fee_amount`) (`FR-04`).
  - Mobile screens: Auth/Onboarding, `Tambah` transaction sheet/screen, `Transaksi` list with search/filter, Account & Category management, `Profil` settings.
- **Dependencies:** Phase 0.

### Phase 2 — Budgets, Reports, Dashboard & Data Export
- **Scope:**
  - PostgreSQL migrations for `budgets`, `recurring_transactions`, and `notification_events`.
  - Deterministic reporting endpoints (`/api/v1/reports/summary`, `/categories`, `/trends`, `/accounts`) and CSV export (`/api/v1/reports/export.csv`) (`FR-08`, `FR-11`).
  - Monthly and per-category budget management with configurable alert thresholds (`FR-07`).
  - Idempotent recurring transaction scheduler (`FR-09`) and budget threshold notifications (`FR-10`).
  - Mobile screens: `Beranda` (Home dashboard with balance privacy toggle, period selector, cash flow summary, spending trend), `Anggaran` (budget progress bars capped visually at 100% with explicit overspent values, status chips).
- **Dependencies:** Phase 1 (requires authoritative transactions, accounts, categories, and transfers).

### Phase 3 — ML Foundation & Assisted Categorization
- **Scope:**
  - PostgreSQL migration for `ml_predictions` and `user_feedback`.
  - Curated Indonesian transaction dataset preparation, train/validation/test split with leak prevention, text normalization pipeline (`services/ml/training/`).
  - Baseline `TF-IDF + LogisticRegression` model training, evaluation report (`Macro-F1`, per-class precision/recall, confusion matrix, confidence calibration) (`services/ml/evaluation/`).
  - Versioned artifact loading in FastAPI and `POST /internal/v1/predict/category` endpoint (`AI-01`).
  - Go `mlclient` integration into transaction creation/suggestion workflow with strict timeout, confidence threshold policy, and fallback to manual selection.
  - Mobile `Tambah` form integration showing clearly labeled auto-suggested categories and one-tap manual override.
- **Dependencies:** Phase 1 (categories & transactions). Can run in parallel with Phase 2.

### Phase 4 — Statistical Insights, Anomaly Detection & Expense Forecasting
- **Scope:**
  - PostgreSQL migration for `insights`.
  - Deterministic evidence aggregation in Go (`services/api/internal/insights`) paired with statistical anomaly detection (`POST /internal/v1/detect/anomalies`, Robust Z-score / historical comparison, `AI-02`) and monthly expense forecasting (`POST /internal/v1/predict/expense-forecast`, historical baseline vs. moving average / exponential smoothing evaluated via MAE/WAPE, `AI-03`).
  - Cold-start guards: suppress forecasts and personalized insights when user transaction history is below the minimum threshold.
  - Mobile `Insight` screen displaying evidence-backed insight cards, comparison periods, forecast estimates with uncertainty ranges, and user relevance feedback (`AI-04`).
- **Dependencies:** Phase 2 (reporting aggregations) and Phase 3 (ML service & evaluation harness).

### Phase 5 — Hardening, Security Audit & Controlled Release
- **Scope:**
  - End-to-end IDOR and security testing, rate limiting verification, log sanitization audit (`SEC-01` to `SEC-05`).
  - Account deletion & data retention workflow verification (`FR-11`).
  - Database backup and restore drill (`15.2`).
  - Performance benchmarking (p95 read latency <= 500 ms target) and Android accessibility/TalkBack/font-scaling QA across light and dark themes (`design.md` Section 14).
- **Dependencies:** Phases 1–4.

---

## 5. Initial Database and API Implementation Strategy

### 5.1 Database Strategy (PostgreSQL)
- **Migration Management:** Sequential, versioned SQL migration pairs (`NNNNNN_name.up.sql` and `NNNNNN_name.down.sql`) stored in `services/api/migrations/`.
- **Monetary Representation:** All monetary columns (`opening_balance`, `amount`, `fee_amount`, budget `amount`) use `BIGINT` representing exact integer minor/base currency units (e.g., whole Rupiah for `IDR`), accompanied by an explicit ISO 4217 `currency` column (`VARCHAR(3)`, default `'IDR'`). IEEE 754 floating-point types (`REAL`, `DOUBLE PRECISION`, `float64`) are strictly forbidden for monetary storage and calculation.
- **Sign Convention:** Stored transaction and transfer amounts are strictly positive (`amount > 0`), enforced via database `CHECK` constraints (`fee_amount >= 0`). Direction is determined deterministically by `transactions.type` (`'income'` vs `'expense'`) or transfer role (`source_account_id` vs `destination_account_id`).
- **Internal Transfer Integrity:** Stored in a dedicated `transfers` table with `CHECK (source_account_id <> destination_account_id)`. Creating, updating, or deleting a transfer runs inside a `SERIALIZABLE` or `READ COMMITTED` row-locked database transaction that validates both accounts belong to the authenticated `user_id`. Principal transfer amounts are excluded from external income/expense reports and budget usage; only `fee_amount > 0` is counted as an external expense.
- **UUID Primary Keys:** Entities use UUIDv4/UUIDv7 primary keys so resource IDs are non-enumerable.

### 5.2 API Implementation Strategy
- **Public API (`/api/v1`):**
  - Standardized JSON success and error envelopes (`{"data": ..., "meta": ...}` and `{"error": {"code": "...", "message": "...", "details": ...}}`).
  - Cursor or limit/offset pagination with bounded `limit` (default 20, max 100) on collection endpoints.
  - `Idempotency-Key` header support on `POST /api/v1/transactions` and `POST /api/v1/transfers` to prevent duplicate financial entries on mobile network retries.
- **Internal ML API (`/internal/v1`):**
  - Bound to the internal Docker network and authenticated via `X-Internal-Service-Token`.
  - Strict timeout enforced by Go `mlclient` (default `1500ms` for interactive category suggestion) with circuit-breaker/fallback behavior so a slow or unavailable ML container never blocks transaction persistence.
  - Every ML response includes `model_version` and structured status (`ready`, `model_not_loaded`, `insufficient_data`).

---

## 6. Testing and Quality Assurance Requirements

1. **Golang (`services/api`):**
   - `gofmt` formatting check and `go vet ./...` static analysis must pass with zero warnings.
   - Unit tests for monetary arithmetic, balance reconciliation, transfer invariants, budget utilization math (including 0-limit and >100% overspent edge cases), input validation, and `mlclient` timeout/fallback handling.
   - Integration tests against PostgreSQL verifying atomic rollbacks on failed transfers and strict `user_id` isolation (IDOR prevention: User A cannot read, mutate, or transfer to/from User B's accounts, categories, budgets, or transactions).
2. **Python (`services/ml`):**
   - Static syntax/type checking and unit/contract tests (`pytest` / `unittest`) covering health endpoints, authentication header checks, Pydantic schema validation, text preprocessing edge cases (empty strings, emojis, noisy merchant strings, typos), and graceful handling when no trained model artifact is loaded.
   - Evaluation gate tests verifying that any trained model artifact produces a reproducible evaluation report on a held-out test set before activation.
3. **Flutter (`apps/mobile`):**
   - `flutter analyze` with `flutter_lints` enabled.
   - Unit tests for IDR currency formatting (`Rp2.840.000`, `+Rp500.000`, `−Rp25.000`, hidden balance masking) and form validators.
   - Widget tests verifying the 5-destination navigation shell, theme token application (light and dark modes), and explicit empty/loading/error/ML-unavailable states.

---

## 7. Security and Financial Data Integrity Requirements

- **Strict Resource Authorization (`SEC-01`):** Every SQL query in `services/api` for user-owned resources must include `WHERE user_id = $1` derived from the verified authentication context—never from an unverified request body or query parameter.
- **Credential & Secret Hygiene (`SEC-01`, `SEC-04`):** Passwords hashed using `argon2id` (or `bcrypt` cost >= 12). Zero secrets committed to Git; `.env` files ignored via `.gitignore`; `.env.example` contains safe placeholders only.
- **Data Minimization for AI (`SEC-03`):** When calling `services/ml`, `services/api` sends only the minimum required features (e.g., normalized description, transaction type, currency, amount) and never transmits user emails, passwords, session tokens, or unrelated account identifiers.
- **Log Redaction (`SEC-02`):** Structured logs (`log/slog` in Go, structured logging in Python) must never log `Authorization` headers, `X-Internal-Service-Token`, passwords, or raw transaction descriptions by default.
- **Deterministic Financial Reconciliation:** Account balances and report summaries must always be reconcilable to `opening_balance + SUM(income) - SUM(expense) - SUM(transfers_out + transfer_fees) + SUM(transfers_in)`.

---

## 8. AI/ML Development and Evaluation Strategy

1. **No Fake Predictions Policy:** Until a model is trained and evaluated on a dataset, `services/ml` must report `model_loaded: false` and return an explicit unavailable/untrained status (`503 Service Unavailable` with error code `MODEL_NOT_READY` or `status: "insufficient_data"`), never hardcoded rules masquerading as ML confidence scores.
2. **Categorization (`AI-01`):**
   - **Baseline Model:** `TfidfVectorizer` (word + character n-grams for Indonesian merchant names/abbreviations) + `LogisticRegression` (with class weights / calibration).
   - **Evaluation Metrics:** Macro-F1 (development target >= 0.80 on representative held-out test set), per-category precision/recall, confusion matrix, and confidence calibration curve.
   - **Confidence Policy:** Only suggestions with calibrated confidence >= configured threshold (default `0.75`) are pre-selected as suggested categories; lower-confidence predictions prompt the user to select manually.
3. **Anomaly Detection (`AI-02`):**
   - **Baseline Approach:** Deterministic historical period comparison + Robust Z-score (median/MAD) per category and transaction size, avoiding false alarms on recurring bills.
   - **Presentation Guardrail:** Must always include numeric evidence and comparison period, and must never label an anomaly as "fraud" or "penipuan".
4. **Monthly Expense Forecasting (`AI-03`):**
   - **Baseline Approach:** Day-of-month run-rate projection and historical moving average, evaluated against actual end-of-month totals using MAE and WAPE.
   - **Cold-Start Guardrail:** Requires minimum historical coverage (default: >= 30 days of history or >= 15 expense transactions across >= 2 weeks) before producing a forecast.

---

## 9. Identified Ambiguities, Technical Risks, and Documented Assumptions

### 9.1 Reconciled Specification Inconsistencies

| Topic | `MonTra_PRD.md` Statement | `design.md` Statement | Documented Resolution |
|---|---|---|---|
| **Primary Bottom Navigation** | Section 10.2 lists 5 tabs: `Home`, `Transactions`, `Budgets`, `Insights`, `Profile`, and states Add Transaction must be easy to find from the dashboard. | Section 4.1 specifies 5 persistent bottom nav items: `Beranda`, `Transaksi`, `Tambah` (center action), `Anggaran`, `Insight`, and explicitly states `Profil` is opened via the top app bar avatar and must **not** be a 6th bottom nav item. | Follow `design.md` for UI navigation layout: 5 bottom destinations (`Beranda`, `Transaksi`, `Tambah`, `Anggaran`, `Insight`) + top app bar avatar opening `Profil`. This satisfies both specifications. |
| **Internal Transfers in Add Flow & Data Model** | FR-03 mentions `income`, `expense`, or `transfer` transaction types, while FR-04, Section 12.1, and Section 13.4 model `transfers` as a dedicated table and `/api/v1/transfers` endpoint. | Section 6.3 (`Tambah`) lists `Pengeluaran` or `Pemasukan` in the basic form fields, while Section 8.5 requires preventing internal transfers from being counted as ordinary income/expense. | Keep `transactions` (`income`, `expense`) and `transfers` as separate database tables and API resources so transfers are never double-counted; in the mobile `Tambah` sheet, support switching between `Pengeluaran`, `Pemasukan`, and `Transfer Antar-Akun`. |
| **Missing HTML Prototype Reference** | Not referenced in PRD. | Line 10 references `montra-android-preview.html`, which is not present in the repository. | `design.md` Sections 5–7 contain the complete, authoritative color hex codes, typography scale, spacing scale, and screen specifications; use `design.md` directly. |
| **PRD File Path** | Prompt mentions `PRD.md` or existing `MonTra_PRD.md`. | Repository contains `docs/MonTra_PRD.md` and `docs/design.md`. | Preserve `docs/MonTra_PRD.md` and `docs/design.md` in `docs/` without moving or deleting them, and reference their exact paths in `README.md` and `AGENTS.md`. |

### 9.2 Documented Technical Assumptions (for Open Decisions in PRD Section 21)

1. **Monetary Precision:** Amounts are stored and transmitted as 64-bit integers (`int64` / PostgreSQL `BIGINT`) in whole Currency Units for `IDR` (since functional IDR transactions do not use sub-rupiah sen in everyday personal finance), alongside an explicit `currency` code (`IDR`).
2. **Authentication Mechanism:** Short-lived JWT access tokens (15 minutes) signed with HMAC-SHA256 (`JWT_SECRET`) paired with database-backed revocable refresh tokens (30 days) for secure logout and session revocation.
3. **Account Balance Reconciliation:** Account balance is defined deterministically as `opening_balance + income_sum - expense_sum - transfer_out_sum - transfer_fee_sum + transfer_in_sum`.
4. **Confidence Threshold:** Default ML categorization auto-suggestion confidence threshold is `0.75`, configurable via environment variable `ML_CATEGORY_CONFIDENCE_THRESHOLD`.
5. **Forecast Cold-Start Threshold:** Minimum 30 calendar days of user history and at least 15 recorded expense transactions before monthly forecasting is enabled; otherwise the API returns `status: "insufficient_data"`.

---

## 10. Explicit Acceptance Criteria by Phase

### Phase 0 — Foundation (Current Iteration)
- [x] `docs/IMPLEMENTATION_PLAN.md` documents architecture, boundaries, phases, risks, assumptions, and acceptance criteria.
- [x] Repository structure is initialized with `services/api`, `services/ml`, `apps/mobile`, `packages/api-spec`, and `infrastructure/`.
- [x] `services/api` compiles cleanly, passes `gofmt`, `go vet`, and `go test ./...`, and exposes `/healthz` and `/api/v1/health`.
- [x] `services/ml` exposes `/internal/v1/health`, enforces contract schemas, does not fabricate ML predictions, and passes automated tests.
- [x] `apps/mobile` establishes the centralized `design.md` light/dark theme tokens, IDR currency formatter, and 5-item navigation shell without hardcoded fake financial records.
- [x] `packages/api-spec/openapi.yaml` defines initial consistent contracts for public and internal endpoints.
- [x] `docker-compose.yml` passes `docker compose config` validation and isolates `ml` on an internal network.
- [x] `.gitignore` and `.env.example` ensure no secrets or local environment files are tracked by Git.

### Phase 1 — Core Financial Ledger & Authentication
- [ ] Database migrations cleanly apply and roll back `users`, `accounts`, `categories`, `transactions`, and `transfers`.
- [ ] Unauthenticated requests to `/api/v1/*` protected routes return `401 Unauthorized`.
- [ ] Cross-user resource access attempts (IDOR tests) return `404 Not Found` or `403 Forbidden`.
- [ ] Internal transfers atomically debit `source_account_id` (plus optional `fee_amount`) and credit `destination_account_id`; simulated mid-transaction failures leave zero partial balance changes.
- [ ] Mobile user can register, log in, create an account, and record/edit/delete an income, expense, or transfer.

### Phase 2 — Budgets, Reports, Dashboard & Data Export
- [ ] Dashboard and report totals match underlying transactions to the exact Rupiah in reconciliation tests; internal transfer principals are excluded from income/expense totals.
- [ ] Budget progress calculations handle 0-limit, <80% (`Dalam batas`), >=threshold (`Mendekati batas`), and >100% (`Melebihi batas`) states accurately, capping visual progress bars at 100% while displaying exact overspent amounts in text.
- [ ] CSV export streams only the authenticated user's filtered transactions.

### Phase 3 — ML Foundation & Assisted Categorization
- [ ] Training and evaluation pipeline runs reproducibly on split datasets with zero test leakage and generates a versioned evaluation artifact.
- [ ] `POST /internal/v1/predict/category` returns `predicted_category`, `confidence`, and `model_version`.
- [ ] Stopping the `ml` container does not prevent creating or editing transactions in `services/api` or `apps/mobile`.

### Phase 4 — Statistical Insights, Anomaly Detection & Expense Forecasting
- [ ] Every generated insight includes `period_start`, `period_end`, comparison period, and verifiable numeric `evidence`.
- [ ] Users with insufficient transaction history receive an explicit calm empty/insufficient-data state rather than fabricated forecasts.
- [ ] Anomaly alerts include measurable historical comparisons and never use fraud/scam accusation terminology.

### Phase 5 — Hardening, Security Audit & Controlled Release
- [ ] All global acceptance criteria in `MonTra_PRD.md` Section 17.5 and visual/usability QA checklists in `design.md` Section 14 pass verification.
