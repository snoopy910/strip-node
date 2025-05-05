# Build stage
FROM golang:1.24-alpine AS builder

# Set GOMODCACHE to use the cache that will be mounted
ENV GOMODCACHE=/go/pkg/mod

WORKDIR /build

RUN apk add build-base curl

# First copy the entire project to handle local module replacements
COPY . .

# Then set up Go module caching
# Create a separate layer for dependency downloads
RUN cd strip-validator && go mod download

# Build the application
RUN cd strip-validator && go build -o strip-validator

# Final stage
FROM alpine:latest

# Install curl for healthcheck
RUN apk add curl

# Create a non-root user for security
RUN adduser -D appuser

WORKDIR /app
# Copy only the binary from the builder stage
COPY --from=builder /build/strip-validator/strip-validator /app/strip-validator

# Install grpc-health-probe
# Find latest version at: https://github.com/grpc-ecosystem/grpc-health-probe/releases
ARG GRPC_HEALTH_PROBE_VERSION=v0.4.3 # Use the latest appropriate version
RUN apk add --no-cache curl ca-certificates && \
    curl -fL https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 -o /usr/local/bin/grpc_health_probe && \
    chmod +x /usr/local/bin/grpc_health_probe && \
    # Clean up (optional, good practice)
    apk del --purge curl && \
    rm -rf /var/cache/apk/*


# Use the non-root user
USER appuser

# Set the entrypoint to the binary
ENTRYPOINT ["/app/strip-validator"]
