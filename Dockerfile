FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY portifyx-goutils /portifyx-goutils
COPY portifyx-api-gateway /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/gateway main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bin/gateway /app/gateway
COPY portifyx-api-gateway/.env.example /app/.env

EXPOSE 8080

CMD ["/app/gateway", "api"]
