.PHONY: api-build api-test lambda-build frontend-build docker-up docker-down docker-build docker-logs lint fmt

COMPOSE ?= docker compose

api-build:
	go build ./src/cmd/api

api-test:
	go test ./...

lambda-build:
	cd src/cmd/lambdas/transaction_processor && GOOS=linux GOARCH=amd64 go build -o bin/transaction_processor

frontend-build:
	cd src/frontend && npm ci && npm run build

docker-up:
	$(COMPOSE) up --build --remove-orphans

docker-down:
	$(COMPOSE) down --remove-orphans

docker-build:
	$(COMPOSE) build

docker-logs:
	$(COMPOSE) logs -f api frontend localstack mongo

fmt:
	gofmt -w ./src/cmd ./src/internal
