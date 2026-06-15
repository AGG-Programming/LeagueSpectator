FROM golang:1.26.4-alpine3.24 AS builder
WORKDIR /app

COPY go.mod ./
COPY ["cmd/main.go", "./cmd/main.go"]

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-validator ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth-validator .
EXPOSE 8081
CMD ["./auth-validator"]