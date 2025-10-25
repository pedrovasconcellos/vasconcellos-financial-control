# Finance Control Platform

Finance Control is a clean-architecture Go backend paired with a React + Material UI frontend for managing personal finances. It integrates with MongoDB for persistence, Amazon S3 and SQS (via LocalStack in local environments), and Cognito for authentication. Asynchronous workloads are handled by an AWS Lambda written in Go that consumes the transactional queue.

## Features

- **Authentication** powered by Amazon Cognito, with a local fallback profile for offline development.
- **Account, category, transaction, budget, goal and reporting APIs** organised strictly around clean architecture boundaries.
- **Async pipeline** where each transaction publishes an SQS message processed by a Lambda that updates budget execution.
- **Asset storage** through S3 for transaction receipts (presigned URLs returned to clients).
- **TypeScript React UI** with Material UI, React Query, and modern UX for dashboards, accounts, transactions, budgets, and goals.
- **Docker-ready** setup with MongoDB and LocalStack, plus bootstrap scripts that provision queues, topics, and Cognito assets.

## Repository layout

```
.
├── cmd/
│   └── api/               # Go HTTP API entry point
├── internal/              # Clean architecture core (entities, use cases, infra)
├── src/
│   └── lambdas/
│       └── transaction_processor/ # AWS Lambda handler
├── frontend/              # React application (Vite + MUI)
├── config/                # Environment configuration templates
├── scripts/               # LocalStack bootstrap scripts
├── Dockerfile             # API build
├── docker-compose.yml     # API + MongoDB + LocalStack + Frontend
└── Makefile               # Handy build/test targets
```

## Getting started

### 1. Prerequisites

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose

### 2. Configure credentials

Copy the template and adjust values as needed (for local use you can keep the defaults):

```bash
cp config/local_credentials.example.yaml config/local_credentials.yaml
```

You can extend `local.authUsers` with additional local login accounts. When running against AWS environments, set `CONFIG_FILE` to the path managed by Secrets Manager that exposes Cognito and database credentials.

### 3. Boot services with Docker

```bash
docker-compose up --build
```

The compose file starts:

- `api`: Go API listening on `http://localhost:8080`
- `frontend`: Vite dev server on `http://localhost:5173`
- `mongo`: MongoDB on `mongodb://localhost:27017`
- `localstack`: Mimics Cognito, S3, and SQS, with `scripts/localstack/00-bootstrap.sh` provisioning required resources automatically.

### 4. Local development without Docker

```bash
# Backend
make api-build
make api-test

# Lambda (binary ready for deployment)
make lambda-build

# Frontend
cd frontend
npm install
npm run dev
```

Set `CONFIG_FILE=config/local_credentials.yaml` before launching the API locally to leverage the offline credentials.

### 5. Running tests and builds

- **Go tests:** `go test ./...`
- **Lambda build:** `GOOS=linux GOARCH=amd64 go build -o bin/transaction_processor ./src/lambdas/transaction_processor`
- **Frontend build:** `npm run build`

### API overview

- `POST /api/v1/auth/login` – Delegates to Cognito or local credentials.
- `GET/POST/PATCH/DELETE /api/v1/accounts`
- `GET/POST/DELETE /api/v1/categories`
- `GET/POST/PATCH /api/v1/transactions`
- `POST /api/v1/transactions/:id/receipt` – Uploads receipt to S3 and returns a presigned URL.
- `GET/POST /api/v1/budgets`
- `GET/POST /api/v1/goals` and `POST /api/v1/goals/:id/progress`
- `GET /api/v1/reports/summary`

Every protected route expects a Bearer token validated via Cognito (or local session when running offline).

### Lambda pipeline

Transactions trigger an SQS message (`finance-transactions-queue`). The Go Lambda consumes the queue, enriches budget metrics, and keeps MongoDB values synchronised. The handler shares the same domain models and repositories as the API, ensuring business invariants stay consistent.

## Environment variables

Key variables honoured by the API:

| Variable | Description |
|----------|-------------|
| `CONFIG_FILE` | Path to a YAML configuration (local or Secrets Manager materialised file). |
| `APP_PORT` (via config file) | HTTP port, defaults to 8080. |
| `AWS_*` fields | Credentials and endpoints for AWS/localstack. |
| `AUTH_MODE` | `cognito` or `local`. |

Refer to `config/local_credentials.example.yaml` for the full schema.

## Frontend integration

The React app looks for `VITE_API_URL` (defaults to `http://localhost:8080/api/v1`) and handles authentication tokens automatically through the shared context provider. Material UI powers a responsive layout, while React Query keeps data fresh across the dashboard sections.

## Contributing

1. Ensure formatting with `gofmt` and `npm run lint` (if you add lint rules) before committing.
2. Keep business code inside the `internal` tree and expose adapters via constructor functions.
3. Extend the documentation in `PROJECT.md` whenever architectural decisions evolve.

Happy hacking! 🚀
