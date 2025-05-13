FROM golang:1.23-alpine AS builder
WORKDIR /cmd/users
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o main ./cmd/users/main.go

FROM alpine:latest
WORKDIR /cmd/users
COPY .env ./
COPY --from=builder /cmd/users/main .
EXPOSE 8088
CMD ["./main"]
