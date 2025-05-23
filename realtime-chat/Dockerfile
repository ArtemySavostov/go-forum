FROM golang:1.23-alpine AS builder
WORKDIR /cmd/chat
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o main ./cmd/chat/main.go

FROM alpine:latest
WORKDIR  /cmd/chat
COPY --from=builder /cmd/chat/main .
EXPOSE 8082
CMD ["./main"]
