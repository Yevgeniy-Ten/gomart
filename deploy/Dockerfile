FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gophermart ./cmd/gophermart

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/gophermart /app/gophermart

CMD ["/app/gophermart"]
