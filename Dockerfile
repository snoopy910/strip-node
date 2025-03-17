FROM golang:1.19-alpine
RUN apk add build-base curl

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build

EXPOSE 8080

ENTRYPOINT ["/app/strip-node"]