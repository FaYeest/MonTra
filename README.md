# MonTra — Personal Finance Tracker with Assistive AI/ML

MonTra (Money Tracker) is an Android personal finance application designed for Indonesian users (`Bahasa Indonesia`, default currency `IDR`). It combines deterministic financial tracking (income, expenses, internal transfers, virtual accounts, budgets, and reports) with explainable, assistive AI/ML capabilities (automatic category suggestions, statistical anomaly detection, and monthly expense forecasting).

## Documentation

- **Product Requirements Document (PRD):** [`docs/MonTra_PRD.md`](./docs/MonTra_PRD.md)
- **Android UI/UX Design Specification:** [`docs/design.md`](./docs/design.md)
- **Architecture & Implementation Plan:** [`docs/IMPLEMENTATION_PLAN.md`](./docs/IMPLEMENTATION_PLAN.md)
- **Agent & Contributor Guidelines:** [`AGENTS.md`](./AGENTS.md)
- **OpenAPI Specification:** [`packages/api-spec/openapi.yaml`](./packages/api-spec/openapi.yaml)

## Repository Structure

```text
montra/
├── AGENTS.md
├── README.md
├── .env.example
├── .gitignore
├── docker-compose.yml
├── docs/
│   ├── MonTra_PRD.md
│   ├── design.md
│   └── IMPLEMENTATION_PLAN.md
├── apps/
│   └── mobile/              # Flutter Android application
├── services/
│   ├── api/                 # Golang authoritative REST API backend
│   └── ml/                  # Python FastAPI internal AI/ML service
├── packages/
│   └── api-spec/            # OpenAPI 3.1 contract definitions
└── infrastructure/
    ├── docker/              # Database bootstrap scripts
    └── proxy/               # Reverse proxy configuration (Caddy)
```

## Architecture Principles

1. **Deterministic Financial Core (`services/api` + PostgreSQL):** Golang owns authentication, authorization, transactions, accounts, internal transfers, budgets, and financial reports. Monetary amounts use exact integer representation (`int64` / `BIGINT`) with explicit currency codes (`IDR`).
2. **Assistive & Isolated AI/ML (`services/ml`):** Python FastAPI handles model training, evaluation, and internal inference (`/internal/v1/*`). It is isolated on an internal network, cannot directly mutate financial records, and never blocks core transaction management if unavailable.
3. **Calm, Evidence-Led Mobile UI (`apps/mobile`):** Flutter implements the design tokens and usability rules from `docs/design.md` with Bahasa Indonesia copy, consistent IDR formatting (`Rp2.840.000`), and explicit loading, empty, error, and ML-unavailable states.

## Local Development Setup

### 1. Environment Configuration

Copy `.env.example` to `.env` in the repository root and customize values if needed for local development:

```powershell
Copy-Item .env.example .env
```

### 2. Run via Docker Compose

When Docker Desktop / Docker Engine is running:

```powershell
# Validate the Compose configuration
docker compose config

# Start PostgreSQL, Go API, Python ML service, and Reverse Proxy
docker compose up --build -d

# Check Go API health (exposed on port 8080, or 8088 via reverse proxy)
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/health
```

> **Note:** `services/ml` is intentionally attached only to the internal Docker network (`internal_net`) and is not published to host ports by default.

### 3. Run Golang Backend Locally (`services/api`)

Requires Go 1.24+ (tested with Go 1.27):

```powershell
# Run static checks and unit tests
go test -v ./...       # Run from services/api
go vet ./...           # Run from services/api

# Start the API server (can run without PostgreSQL/ML in degraded health check mode or with local services)
go run ./cmd/api
```

### 4. Run Python AI/ML Service Locally (`services/ml`)

Requires Python 3.12+:

```powershell
# Create a virtual environment in services/ml/.venv and install dependencies
py -3.14 -m venv services/ml/.venv
.\services\ml\.venv\Scripts\pip install -e "services/ml[dev]"

# Run unit and API contract tests
.\services\ml\.venv\Scripts\pytest services/ml/tests

# Start the internal FastAPI server on localhost:8000 for local development
.\services\ml\.venv\Scripts\uvicorn app.main:app --app-dir services/ml --host 127.0.0.1 --port 8000
```

### 5. Run Flutter Mobile App (`apps/mobile`)

Requires Flutter SDK (3.24+ / Dart 3.5+):

```powershell
# Fetch dependencies, analyze, and run widget/unit tests (from apps/mobile)
flutter pub get
flutter analyze
flutter test
```
