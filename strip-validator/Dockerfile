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

# Use the non-root user
USER appuser

# Set the entrypoint to the binary
ENTRYPOINT ["/app/strip-validator"]
