# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.24 AS builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

WORKDIR /workspace

# Instala dependências do módulo principal
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copia todo o código
COPY . .

# Compila o binário com strip de símbolos para reduzir tamanho
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /workspace/bin/api ./cmd/api

FROM --platform=$TARGETPLATFORM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /workspace/bin/api /app/api
COPY config /app/config

ENV CONFIG_FILE=/app/config/local_credentials.yaml
EXPOSE 8080

ENTRYPOINT ["/app/api"]
