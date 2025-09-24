FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/auth/server/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

CMD {"./server"}