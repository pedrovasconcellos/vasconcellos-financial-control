# Financial Control Platform

> **For AI agents and developers:**
> - Agent guidelines: `AGENTS.md` (includes documentation map)
> - Architecture decisions: `PROJECT.md`
> - Known issues: `BUGS_AND_LIMITATIONS.md`
> - Planned improvements: `IMPROVEMENTS.md`

## Overview

Financial Control is a personal finance platform built on top of a clean-architecture Go backend and a React (Vite + Material UI) frontend. The system persists data in MongoDB, integrates with Amazon S3 and SQS (optionally via LocalStack for local development), relies on Amazon Cognito for authentication, and offloads budget recalculations to a Go-based AWS Lambda that consumes the transaction queue.

## Architecture Highlights

- **Clean architecture backend** – domain entities, repositories, and use cases live under `src/internal/`; HTTP adapters and infrastructure concerns stay in their own packages.
- **Type-safe frontend** – React + TypeScript with Material UI, React Router, and React Query for state and data access.
- **Asynchronous pipeline** – every recorded transaction emits an SQS message; the Lambda updates budget execution totals.
- **Secure receipts** – transaction receipts are encrypted with AES-256 and stored in S3; presigned URLs are returned to the UI.
- **Configurable environments** – configuration is composed from defaults, a YAML file referenced via `CONFIG_FILE`, and environment variables.

## Directory Structure

```
.
├── src/
│   ├── cmd/                         # Application entrypoints
│   │   ├── api/                     # HTTP API entrypoint
│   │   └── lambdas/
│   │       └── transaction_processor/   # AWS Lambda handler
│   ├── internal/                    # Domain entities, use cases, and adapters
│   ├── frontend/                    # React application (Vite + Material UI)
│   └── configs/                     # Configuration templates
├── scripts/                         # Utilities (LocalStack bootstrap, migrations, seeds)
├── infra/terraform/                # Optional AWS provisioning stack
├── docker/                          # Dockerfiles
├── docker-compose.yml
└── Makefile
```

## Identifier Strategy

All collections now use UUID strings (`uuid.NewString()` in Go).

## Requirements

| Tool | Version |
|------|---------|
| Go  | 1.24 or newer |
| Node.js | 20.19 or newer (required by Vite 7) |
| npm | ships with Node 20 |
| Docker & Docker Compose | latest stable |
| AWS CLI | configured for homolog/production deployments |
| AES-256 key | Base64-encoded 32-byte key for receipt encryption (`security.encryptionKey`) |

> **Secrets**: Outside local development, store sensitive values (encryption key, database credentials, AWS keys) in AWS Secrets Manager or Systems Manager Parameter Store and inject them into the runtime. Never version encryption keys.

## Local Development

1. **Configure credentials**
   ```bash
   cp src/configs/local_credentials.example.yaml src/configs/local_credentials.yaml
   ```
   Generate an encryption key (`openssl rand -base64 32`) and fill `security.encryptionKey`. The sample file already points to the services created by Docker Compose.

2. **Start the stack**
   ```bash
   docker compose up --build
   ```
   This boots:
   - API on `http://localhost:8080` (`/api/v1/health` for readiness)
   - Frontend on `http://localhost:5173`
   - MongoDB on `mongodb://localhost:27017`
   - LocalStack (S3, SQS, Cognito) configured by `scripts/localstack/00-bootstrap.sh` (creates `financial-transactions-queue` with a dead-letter queue `financial-transactions-dlq` and a redrive policy of 5 attempts)
   - Transaction lambda worker (`transaction-lambda` service) polling the queue continuously with `LAMBDA_LOCAL=true`

3. **Run services manually (optional)**
   ```bash
   # API
   export CONFIG_FILE=src/configs/local_credentials.yaml
   go run ./src/cmd/api

   # Lambda build (for local tests)
   GOOS=linux GOARCH=amd64 go build -o bin/transaction_processor ./src/cmd/lambdas/transaction_processor

   # Frontend with hot reload
   cd src/frontend
   npm install
   npm run dev
   ```

4. **Seed data**
   ```bash
   # Complete seed script (creates users, accounts, categories, transactions, budgets, and goals)
   mongosh financial-control scripts/seed_complete.js
   ```
   Adjust the `docker exec ... mongosh` command if you are running MongoDB inside the compose stack (container name defaults to `financial-control-mongo-1`).

5. **Tests**
   ```bash
   go test ./...
   npm run build        # Type-checks the frontend
   make lambda-build    # Produces the Lambda artifact
   ```

## Data Seeding Reference

| Script | Purpose |
|--------|---------|
| `scripts/seed_complete.js` | Complete seed script that creates two users (vasconcellos and teste), their accounts, categories, transactions, budgets, and goals for testing and development. |

All scripts are written for `mongosh`; pipe them through `mongosh <database> < script.js` or use `docker exec` when running MongoDB in Docker.

## Deployment (Homolog / Production)

1. **Configuration**
   - Provide a YAML file via `CONFIG_FILE` (mounted secret or rendered during deploy) and complement with environment variables for sensitive values.
   - Set `auth.mode=cognito`, provide the real AWS region, Cognito Client ID, S3 bucket, and queue URL.

2. **Build and publish the backend**
   ```bash
   docker build -t financial-api:latest .
   aws ecr create-repository --repository-name financial-api --region <region>    # one-time
   docker tag financial-api:latest <ACCOUNT_ID>.dkr.ecr.<region>.amazonaws.com/financial-api:latest
   docker push <ACCOUNT_ID>.dkr.ecr.<region>.amazonaws.com/financial-api:latest
   ```
   Deploy the image to ECS Fargate or EC2 (Systemd). Provide:
   ```
   CONFIG_FILE=/app/src/configs/config.yaml
   APP_ENVIRONMENT=homolog|production
   AWS_REGION=<region>
   ```

3. **Lambda deployment**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bootstrap ./src/cmd/lambdas/transaction_processor
   zip lambda.zip bootstrap
   aws lambda update-function-code \
     --function-name financial-transaction-processor \
     --zip-file fileb://lambda.zip
   ```
   Configure the same environment variables (`CONFIG_FILE`, `AWS_REGION`, `AUTH_MODE`) and connect the SQS queue as a trigger. Enable a DLQ for resilience.

4. **Frontend**
   ```bash
   cd src/frontend
   npm install
   npm run build
   ```
Deploy the contents of `src/frontend/dist` to S3 + CloudFront or another static hosting solution. Set `VITE_API_URL` (e.g., `https://api.company.com/api/v1`) in the environment prior to building.

## API Endpoints

- `POST /api/v1/auth/login`
- `GET/POST/PATCH/DELETE /api/v1/accounts`
- `GET/POST/DELETE /api/v1/categories`
- `GET/POST/PATCH /api/v1/transactions`
- `POST /api/v1/transactions/:id/receipt`
- `GET/POST /api/v1/budgets`
- `GET/POST /api/v1/goals`
- `POST /api/v1/goals/:id/progress`
- `GET /api/v1/reports/summary`

`GET` endpoints for accounts, transactions, budgets, and goals accept optional `limit` and `offset` query parameters (`limit` defaults to 100, capped at 200; `offset` defaults to 0) to support pagination on large datasets.

### Common Environment Variables

| Variable | Notes |
|----------|-------|
| `CONFIG_FILE` | Path to the YAML config. |
| `APP_ENVIRONMENT` | `development`, `homolog`, or `production` (used for logging/metrics). |
| `AUTH_MODE` | `local` for dev; `cognito` in managed environments. |
| `security.encryptionKey` | Base64-encoded AES-256 key (must come from a secret in homolog/production). |

Consult `src/configs/local_credentials.example.yaml` for the full schema.

## Terraform Automation (Optional)

The module in `infra/terraform` provisions AWS App Runner for the API, DocumentDB Serverless for MongoDB, S3, SQS (with DLQ), and Cognito. Quick start:

```bash
cd infra/terraform
cp terraform.tfvars.example terraform.tfvars   # customise environment, image URI, budgets, etc.

terraform init
terraform plan
terraform apply
```

Outputs include the App Runner URL, DocumentDB endpoint, Cognito identifiers, and S3 bucket name. Review costs before applying (approx. USD 45-50/month for the default sizing). See `infra/terraform/README_INFRA.md` for detailed infrastructure documentation.

## Useful Make Targets

```bash
make api-build         # go build ./src/cmd/api
make api-test          # go test ./...
make lambda-build      # build Lambda binary (linux/amd64)
make frontend-build    # npm install && npm run build
make docker-up         # start all services with docker-compose
make docker-down       # stop all services
make docker-logs       # view logs from all services
make fmt               # format Go code
```

## Conventions & Further Reading

- Keep domain logic inside `src/internal/usecase` and `src/internal/domain`; adapters should remain thin and testable.
- Comments should be in Portuguese when necessary for context; code remains in English.
- Document architectural or configuration changes in `PROJECT.md`.
- Before opening a pull request, run `go test ./...`, `make lambda-build`, and `npm run build`.

For a deeper dive into design decisions and open technical work, refer to:

- `PROJECT.md` – architectural decisions and conventions.
- `BUGS_AND_LIMITATIONS.md` – known issues and technical debt.
- `IMPROVEMENTS.md` – backlog of enhancements.

Happy hacking! 🚀
