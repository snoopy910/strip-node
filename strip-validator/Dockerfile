# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /build

RUN apk add build-base curl

COPY . .

# RUN go mod vendor
RUN go build -o strip-validator

# Final stage
FROM alpine:latest

# Create a non-root user for security
RUN adduser -D appuser

WORKDIR /app
# Copy only the binary from the builder stage
COPY --from=builder /build/strip-validator /app/strip-validator

RUN apk add curl

# Use the non-root user
USER appuser

# Set the entrypoint to the binary
ENTRYPOINT ["/app/strip-validator"]
