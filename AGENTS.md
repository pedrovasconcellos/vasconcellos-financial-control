# Agent Guidelines

This repository hosts a full-stack personal finance platform. When interacting with it:

- **Architecture:** The Go backend follows clean architecture. Domain types, repositories, and use cases live under `internal/`. Adapters (HTTP handlers, MongoDB repositories, AWS clients) must remain thin and respect dependency direction (outer layers depend on inner ones only).
- **Entrypoints:**
  - API bootstraps from `src/api/cmd/api/main.go`.
  - AWS Lambda handler is under `src/lambdas/transaction_processor/main.go` and shares the same domain packages.
  - React application lives in `frontend/` (Vite + Material UI).
- **Configuration:** Prefer reading configuration via `internal/config.LoadConfig()` which merges environment variables and optional `CONFIG_FILE` YAML. Local development uses `config/local_credentials.yaml` (gitignored) with `auth.mode=local`.
- **AWS integrations:** Interact with AWS through the thin wrappers located in `internal/infrastructure/aws`. They already handle LocalStack endpoints. For Cognito, rely on the auth providers under `internal/infrastructure/auth`.
- **Storage:** MongoDB repositories reside in `internal/infrastructure/mongodb`. When introducing new collections, add indexes in the repository constructors.
- **Async processing:** Publishing to queues happens through the `port.QueuePublisher` interface. The lambda consumes from the same queue and updates budgets. Any new asynchronous feature should follow the same contract.
- **Testing:** Run `go test ./...` for backend code and `npm run build` for the frontend (type-checking). Add unit tests close to the package being tested.
- **Docker/localstack:** `docker-compose.yml` provisions MongoDB, LocalStack, API, and frontend. LocalStack bootstrapping scripts live in `scripts/localstack`.
- **Coding standards:**
  - Go code must be formatted with `gofmt` and keep comments (when necessary) in Portuguese, while code remains in English.
  - React components use functional style, Material UI theming, and React Query for data fetching.
  - Prefer dependency injection via constructors and keep global state isolated to the `providers/` layer on the frontend.

Before shipping changes, run the Go tests, build the lambda (`make lambda-build`), and `npm run build` to ensure the UI type-checks. Document architectural or configuration changes in `PROJECT.md`.
