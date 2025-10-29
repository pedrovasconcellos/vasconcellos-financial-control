# Agent Guidelines

## Purpose

This file (`AGENTS.md`) provides essential guidelines and context for AI agents working with this codebase. It serves as the primary reference for understanding the project architecture, conventions, and development practices. Always refer to this file when making changes or additions to ensure consistency with the project's standards.

## Repository Overview

This repository hosts a full-stack personal financial platform. When interacting with it:

- **Architecture:** The Go backend follows clean architecture. Domain types, repositories, and use cases live under `src/internal/`. Adapters (HTTP handlers, MongoDB repositories, AWS clients) must remain thin and respect dependency direction (outer layers depend on inner ones only).
- **Entrypoints:**
  - API bootstraps from `src/cmd/api/main.go`.
  - AWS Lambda handler is under `src/cmd/lambdas/transaction_processor/main.go` and shares the same domain packages.
  - React application lives in `src/frontend/` (Vite + Material UI).
- **Configuration:** Prefer reading configuration via `src/internal/config.LoadConfig()` which merges environment variables and optional `CONFIG_FILE` YAML. Local development uses `src/configs/local_credentials.yaml` (gitignored) with `auth.mode=local`. **Never commit sensitive values** (encryption keys, AWS credentials); use AWS Secrets Manager in homolog/production.
- **AWS integrations:** Interact with AWS through the thin wrappers located in `src/internal/infrastructure/aws`. They already handle LocalStack endpoints. For Cognito, rely on the auth providers under `src/internal/infrastructure/auth`.
- **Storage:** MongoDB repositories reside in `src/internal/infrastructure/mongodb`. When introducing new collections, add indexes in the repository constructors. All collections use **UUID strings** as identifiers (`uuid.NewString()` in Go). The `processed_transactions` collection ensures idempotency in async transaction processing.
- **Async processing:** Publishing to queues happens through the `port.QueuePublisher` interface. The lambda consumes from the same queue and updates budgets. **Idempotency is critical**: the `processed_transactions` collection prevents duplicate processing. Any new asynchronous feature should follow the same contract and ensure idempotency.
- **Testing:** Run `go test ./...` for backend code and `npm run build` for the frontend (type-checking). Add unit tests close to the package being tested.
- **Docker/localstack:** `docker-compose.yml` provisions MongoDB, LocalStack, API, and frontend. Dockerfiles are organized in `docker/` directory (e.g., `Dockerfile.api`, `Dockerfile.frontend`, `Dockerfile.lambda`). LocalStack bootstrapping scripts live in `scripts/localstack` and automatically create S3 buckets, SQS queues with DLQ, and Cognito resources. **Note:** LocalStack doesn't send real emails; use `auth.mode=local` for development.
- **Coding standards:**
  - Go code must be formatted with `gofmt` and keep comments (when necessary) in Portuguese, while code remains in English.
  - React components use functional style, Material UI theming, and React Query for data fetching.
  - Prefer dependency injection via constructors and keep global state isolated to the `providers/` layer on the frontend.
- **Documentation:** It is important to keep `README.md` updated with important project information. Whenever documenting the project in `README.md` (English), automatically copy and translate all content to `README_PTBR.md` (Portuguese), maintaining the same structure and formatting. Document architectural or configuration changes in `PROJECT.md`. For known issues and improvements, check `BUGS_AND_LIMITATIONS.md` and `IMPROVEMENTS.md`. **All Terraform infrastructure changes must be documented in `infra/terraform/README_INFRA.md`**.

## Requirements

| Tool | Version |
|------|---------|
| Go | 1.24 or newer |
| Node.js | 20.19 or newer (required by Vite 7) |
| npm | ships with Node 20 |
| Docker & Docker Compose | latest stable |
| AWS CLI | configured for homolog/production deployments |

## Quick Setup

1. **Start all services:**
   ```bash
   make docker-up
   ```

2. **Verify services:**
   - API: http://localhost:8080/api/v1/health
   - Frontend: http://localhost:5173
   - MongoDB: mongodb://localhost:27017

## Common Commands

Use these Makefile shortcuts for common operations:
- `make api-build` - Build the API binary
- `make api-test` - Run all Go tests
- `make lambda-build` - Build Lambda binary (linux/amd64)
- `make frontend-build` - Build frontend with npm
- `make docker-up` - Start all services with docker-compose
- `make docker-down` - Stop all services
- `make docker-logs` - View logs from all services
- `make fmt` - Format Go code

## Environment Variables

Key environment variables:
- `CONFIG_FILE` - Path to YAML config file (defaults to `src/configs/local_credentials.yaml`)
- `APP_ENVIRONMENT` - `development`, `homolog`, or `production`
- `AWS_REGION` - AWS region for services (default: `us-east-1`)
- `AUTH_MODE` - `local` for dev, `cognito` for production
- `VITE_API_URL` - Frontend API URL (must be set before `npm run build`)

## Deployment Scripts

- **Frontend deployment:** Use `scripts/deploy-frontend.sh <bucket-name> <api-url>` for S3+CloudFront deployment. The script builds, uploads, and invalidates CloudFront cache automatically.
- **Lambda deployment:** Build with `make lambda-build`, then zip and upload to AWS Lambda.

## Service URLs & Ports

When running locally with `make docker-up`:
- API: http://localhost:8080
- Frontend: http://localhost:5173
- MongoDB: mongodb://localhost:27017
- LocalStack: http://localhost:4566

Health check endpoint: `GET http://localhost:8080/api/v1/health`

## Directory Structure

Key directories:
- `src/internal/domain/entity/` - Business entities (Transaction, Account, Budget, etc.)
- `src/internal/usecase/` - Business logic orchestration
- `src/internal/infrastructure/mongodb/` - MongoDB repository implementations
- `src/internal/adapters/http/handler/` - HTTP request handlers
- `src/internal/adapters/http/middleware/` - HTTP middleware (auth, logging)
- `src/internal/infrastructure/aws/` - AWS service wrappers (S3, SQS)
- `src/frontend/src/pages/` - React page components
- `src/frontend/src/services/` - API client services

## Common Pitfalls & Troubleshooting

- **UUID vs ObjectId**: Always use `uuid.NewString()` for new IDs.
- **Idempotency**: When working with async processing, always check `processed_transactions` collection before processing to avoid duplicates.
- **LocalStack limitations**: LocalStack doesn't send real emails. Always use `auth.mode=local` for local development.
- **CORS in production**: Never use `AllowOrigins: ["*"]` outside development. Check `BUGS_AND_LIMITATIONS.md` for known security issues.
- **Transaction consistency**: Account balance updates and transaction recording should be in the same transaction. See known bugs in `BUGS_AND_LIMITATIONS.md`.

## Known Limitations

Refer to `BUGS_AND_LIMITATIONS.md` for detailed list, but key points:
- MongoDB transactions not fully implemented (balance + transaction may be inconsistent on crash)
- Frontend token refresh not implemented (expires without renewal)
- CORS allows all origins in some configurations (security risk)
- No tracing/metrics (difficult to debug production issues)

## Code Examples

When creating a new repository:
```go
// src/internal/infrastructure/mongodb/new_repository.go
func NewExampleRepository(collection *mongo.Collection) *ExampleRepository {
    // Always create indexes
    collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}},
    })
    return &ExampleRepository{collection: collection}
}
```

When adding a new use case:
```go
// src/internal/usecase/example_usecase.go
// 1. Define repository interface in src/internal/domain/repository/
// 2. Implement in src/internal/infrastructure/mongodb/
// 3. Inject via constructor in src/internal/usecase/
```

## Workflow for Adding New Features

1. **New domain entity:**
   - Create in `src/internal/domain/entity/`
   - Define repository interface in `src/internal/domain/repository/`
   - Implement MongoDB repository in `src/internal/infrastructure/mongodb/` (with indexes!)
   - Create use case in `src/internal/usecase/`
   - Add HTTP handler in `src/internal/adapters/http/handler/`
   - Register route in `src/internal/adapters/http/router.go`

2. **New frontend page:**
   - Create component in `src/frontend/src/pages/`
   - Add route in `src/frontend/src/App.tsx`
   - Create API service method in `src/frontend/src/services/api.ts`

## Testing & Debugging

- **Run specific test:** `go test -v ./src/internal/usecase -run TestAccountUseCase`
- **Run tests with coverage:** `go test -cover ./...`
- **Debug API locally:** Set `CONFIG_FILE=src/configs/local_credentials.yaml` and run `go run ./src/cmd/api`
- **Frontend hot reload:** `cd src/frontend && npm run dev`
- **View logs:** `make docker-logs` or `docker compose logs -f api`
- **Check service health:** `curl http://localhost:8080/api/v1/health`

## Commit & Pull Request Guidelines

- Run all checks before committing: `go test ./...`, `make lambda-build`, `npm run build`
- Commit messages should be clear and descriptive
- Document architectural changes in `PROJECT.md`
- Document infrastructure changes in `infra/terraform/README_INFRA.md`
- Update `README.md` and `README_PTBR.md` simultaneously when documenting

## Quick Reference

- **API Base URL:** http://localhost:8080/api/v1
- **Main endpoints:** `/accounts`, `/transactions`, `/budgets`, `/goals`, `/categories`, `/reports/summary`
- **Authentication:** POST `/auth/login` (local mode: see `local_credentials.yaml` for test users)
- **Database:** MongoDB collection names match entity names (lowercase): `accounts`, `transactions`, `budgets`, `goals`, `categories`, `users`, `processed_transactions`

Before shipping changes, **always run unit tests** after any code changes using `go test ./...`, build the lambda (`make lambda-build`), and `npm run build` to ensure the UI type-checks.
