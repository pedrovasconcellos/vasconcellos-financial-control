# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.24 AS builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

WORKDIR /workspace

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY internal ./internal
COPY cmd ./cmd

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /workspace/bin/api ./cmd/api

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /workspace/bin/healthcheck ./cmd/tools/healthcheck

FROM --platform=$TARGETPLATFORM cgr.dev/chainguard/static:latest AS runtime
WORKDIR /app
COPY --from=builder /workspace/bin/api /app/api
COPY --from=builder /workspace/bin/healthcheck /app/healthcheck
COPY config /app/config
COPY certs /app/certs

ENV CONFIG_FILE=/app/config/local_credentials.yaml
EXPOSE 8080

ENTRYPOINT ["/app/api"]
