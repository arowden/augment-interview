# Stage 1: Build
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build with stripped debug info
ARG VERSION=dev
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -o /app/server ./cmd/server

# Stage 2: Minimal secure runtime
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary with correct ownership
COPY --from=builder --chown=nonroot:nonroot /app/server /server

# Note: migrations are embedded in binary, not copied
USER nonroot:nonroot
WORKDIR /home/nonroot

EXPOSE 8080

# Healthcheck compatible with distroless (no shell)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD ["/server", "--health-check"]

ENTRYPOINT ["/server"]
