FROM golang:1.23-alpine AS builder
WORKDIR /cmd/gateway
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o main ./cmd/gateway/main.go

FROM alpine:latest
WORKDIR  /cmd/gateway
COPY .env ./
COPY --from=builder /cmd/gateway/main .
EXPOSE 8000
CMD ["./main"]
