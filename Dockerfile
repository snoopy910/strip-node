# Build stage
FROM golang:1.19-alpine AS builder

# Set GOMODCACHE to use the cache that will be mounted
ENV GOMODCACHE=/go/pkg/mod

WORKDIR /build

RUN apk add build-base curl

# Copy Go module files first for better layer caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN go build -o strip-node

# Final stage
FROM alpine:latest

# Install curl for healthcheck
RUN apk add curl

# Create a non-root user for security
RUN adduser -D appuser

WORKDIR /app

# Create necessary directories
RUN mkdir -p /app/static-bootnode && \
    chown -R appuser:appuser /app

# Copy only the binary from the builder stage
COPY --from=builder /build/strip-node /app/strip-node

# Copy the static bootnode key file to ensure consistent peer ID
COPY bootnode/static-bootnode/ /app/static-bootnode/

# Use the non-root user
USER appuser

# Expose necessary port
EXPOSE 8080

# Set the entrypoint to the binary
ENTRYPOINT ["/app/strip-node"]