# Build stage
FROM golang:1.19-alpine AS builder
RUN apk add build-base curl

WORKDIR /build

COPY . .

RUN go build -o bootnode

# Final stage
FROM alpine:latest

RUN apk add curl

COPY --from=builder /build/bootnode /app/bootnode
COPY --from=builder /build/static-bootnode/bootnode.bin /static-bootnode/bootnode.bin

ENTRYPOINT ["/app/bootnode"]