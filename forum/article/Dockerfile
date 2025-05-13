FROM golang:1.23-alpine AS builder
WORKDIR /cmd/articles
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o main ./cmd/articles/main.go
RUN ls -l  


FROM alpine:latest
WORKDIR /cmd/articles
COPY .env ./
COPY --from=builder /cmd/articles/main .
EXPOSE 8080
CMD ["./main"]
