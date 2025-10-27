# Agent Guidelines

## Purpose

This file (`AGENTS.md`) provides essential guidelines and context for AI agents working with this codebase. It serves as the primary reference for understanding the project architecture, conventions, and development practices. Always refer to this file when making changes or additions to ensure consistency with the project's standards.

## Repository Overview

This repository hosts a full-stack personal financial platform. When interacting with it:

- **Architecture:** The Go backend follows clean architecture. Domain types, repositories, and use cases live under `internal/`. Adapters (HTTP handlers, MongoDB repositories, AWS clients) must remain thin and respect dependency direction (outer layers depend on inner ones only).
- **Entrypoints:**
  - API bootstraps from `cmd/api/main.go`.
  - AWS Lambda handler is under `cmd/lambdas/transaction_processor/main.go` and shares the same domain packages.
  - React application lives in `frontend/` (Vite + Material UI).
- **Configuration:** Prefer reading configuration via `internal/config.LoadConfig()` which merges environment variables and optional `CONFIG_FILE` YAML. Local development uses `config/local_credentials.yaml` (gitignored) with `auth.mode=local`. **Never commit sensitive values** (encryption keys, AWS credentials); use AWS Secrets Manager in homolog/production.
- **AWS integrations:** Interact with AWS through the thin wrappers located in `internal/infrastructure/aws`. They already handle LocalStack endpoints. For Cognito, rely on the auth providers under `internal/infrastructure/auth`.
- **Storage:** MongoDB repositories reside in `internal/infrastructure/mongodb`. When introducing new collections, add indexes in the repository constructors. All collections use **UUID strings** as identifiers (`uuid.NewString()` in Go). Legacy documents with `ObjectId` must be migrated using `scripts/convert_objectids_to_uuid.js`. The `processed_transactions` collection ensures idempotency in async transaction processing.
- **Async processing:** Publishing to queues happens through the `port.QueuePublisher` interface. The lambda consumes from the same queue and updates budgets. **Idempotency is critical**: the `processed_transactions` collection prevents duplicate processing. Any new asynchronous feature should follow the same contract and ensure idempotency.
- **Testing:** Run `go test ./...` for backend code and `npm run build` for the frontend (type-checking). Add unit tests close to the package being tested.
- **Docker/localstack:** `docker-compose.yml` provisions MongoDB, LocalStack, API, and frontend. Dockerfiles are organized in `docker/` directory (e.g., `Dockerfile.api`, `Dockerfile.frontend`, `Dockerfile.lambda`). LocalStack bootstrapping scripts live in `scripts/localstack` and automatically create S3 buckets, SQS queues with DLQ, and Cognito resources. **Note:** LocalStack doesn't send real emails; use `auth.mode=local` for development.
- **Coding standards:**
  - Go code must be formatted with `gofmt` and keep comments (when necessary) in Portuguese, while code remains in English.
  - React components use functional style, Material UI theming, and React Query for data fetching.
  - Prefer dependency injection via constructors and keep global state isolated to the `providers/` layer on the frontend.
- **Documentation:** It is important to keep `README.md` updated with important project information. Whenever documenting the project in `README.md` (English), automatically copy and translate all content to `README_PTBR.md` (Portuguese), maintaining the same structure and formatting. Document architectural or configuration changes in `PROJECT.md`. For known issues and improvements, check `BUGS_AND_LIMITATIONS.md` and `IMPROVEMENTS.md`. **All Terraform infrastructure changes must be documented in `infra/terraform/README_INFRA.md`**.

Before shipping changes, **always run unit tests** after any code changes using `go test ./...`, build the lambda (`make lambda-build`), and `npm run build` to ensure the UI type-checks.
