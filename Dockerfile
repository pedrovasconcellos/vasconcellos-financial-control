# syntax=docker/dockerfile:1.4

FROM golang:1.21 AS builder
WORKDIR /app

COPY go.work go.work
COPY src/api/go.mod src/api/go.sum ./src/api/
COPY src/lambdas/transaction_processor/go.mod src/lambdas/transaction_processor/go.mod

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd src/api && go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd src/api && go build -o /app/bin/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/bin/api /app/api
COPY config /app/config

ENV CONFIG_FILE=/app/config/local_credentials.yaml
EXPOSE 8080

ENTRYPOINT ["/app/api"]
