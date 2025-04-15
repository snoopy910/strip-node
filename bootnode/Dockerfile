# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /build
RUN apk add build-base curl

ENV ENV=${ENV:-development}

# Copy only necessary files
COPY bootnode.go .
COPY go.mod go.sum ./
COPY /static-bootnode/${ENV} .

RUN go build -o bootnode

# Final stage
FROM alpine:latest

ENV LISTEN_HOST=${LISTEN_HOST:-0.0.0.0}
ENV PORT=${PORT:-30303}

WORKDIR /app
RUN apk add curl

COPY --from=builder /build/bootnode /app/
COPY --from=builder /build/bootnode.bin /app/bootnode.bin

ENTRYPOINT ["/app/bootnode"]
