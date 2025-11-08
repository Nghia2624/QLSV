# syntax=docker/dockerfile:1
FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o query-service ./cmd/query-service

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/query-service .
RUN apk --no-cache add ca-certificates curl
EXPOSE 8081
CMD ["./query-service"] 