FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/migrations /app/migrations

RUN chmod +x /app/api
EXPOSE 8080

CMD ["./api"]