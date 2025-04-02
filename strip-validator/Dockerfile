# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /build

RUN apk add build-base curl

COPY . .

# Run go mod vendor and build
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
