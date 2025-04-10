# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /build
RUN apk add build-base curl

# Copy only necessary files
COPY bootnode.go .
COPY go.mod go.sum ./
COPY static-bootnode ./static-bootnode

RUN go build -o bootnode

# Final stage
FROM alpine:latest

WORKDIR /app
RUN apk add curl

COPY --from=builder /build/bootnode /app/
COPY --from=builder /build/static-bootnode /static-bootnode

# Set environment variables
ENV LISTEN_HOST=0.0.0.0
ENV PORT=30303
ENV KEY_PATH=/static-bootnode

ENTRYPOINT ["/app/bootnode"]
