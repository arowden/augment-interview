# Stage 1: Build
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Cache dependencies
COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY backend/ .

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

# Note: ECS uses ALB target group health checks, not Docker HEALTHCHECK

ENTRYPOINT ["/server"]
