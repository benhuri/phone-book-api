FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o phone-book-api ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/phone-book-api .

CMD ["./phone-book-api"]