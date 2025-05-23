FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go test -c -o integration.test

FROM alpine:latest

RUN apk add --no-cache mysql-client

WORKDIR /app

COPY --from=builder /app/integration.test .

COPY wait-for-mysql.sh /usr/local/bin/wait-for-mysql.sh
RUN chmod +x /usr/local/bin/wait-for-mysql.sh

#ENTRYPOINT ["/usr/local/bin/wait-for-mysql.sh"]
CMD ["./integration.test", "-test.v"]