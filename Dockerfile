FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /phone-book-api cmd/main.go

EXPOSE 8080

CMD ["/phone-book-api"]