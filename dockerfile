FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY cmd/ ./cmd/
RUN go build -o main ./cmd/
CMD ["./main"]