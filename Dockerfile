# Multi-stage build for Testing Studio
# Based on TMS datagen Dockerfile patterns

# Stage 1: Builder base - Download dependencies
FROM australia-southeast1-docker.pkg.dev/anz-x-artifacts-prod/base-images/go:1.25-dev AS builder-base

WORKDIR /workdir
RUN mkdir /builddir

# Configure Go proxy and checksum verification (following TMS pattern)
ENV GONOSUMDB="github.com/anzx/*,github.service.anz/*"
ENV GOPROXY=https://platform-gomodproxy.services.x.gcp.anz/,direct

# Copy go modules
COPY go.mod go.sum ./

# Download dependencies with caching
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x

# Stage 2: Builder - Build the application
FROM builder-base AS builder

ARG entrypoint=cmd/server
WORKDIR /workdir

# Copy source code
COPY . .

# Build the application with build cache
# Note: CGO_ENABLED=1 required for confluent-kafka-go (librdkafka C bindings)
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /builddir/testing-studio ./$entrypoint

# Stage 3: Runtime - Debian slim with glibc for CGO support
FROM debian:bookworm-slim AS runtime

LABEL ci_group="Testing-Studio"
LABEL ci_name="CloudEvents-Explorer"

# Install runtime dependencies for CGO (required by confluent-kafka-go)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1000 appuser && \
    useradd -r -u 1000 -g appuser appuser

WORKDIR /app

COPY --from=builder /builddir/testing-studio /app/testing-studio

# Change ownership
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8888

ENTRYPOINT ["/app/testing-studio"]
