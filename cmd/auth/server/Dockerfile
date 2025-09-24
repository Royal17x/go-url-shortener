FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/auth/server/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server /app/server

RUN useradd --system --uid 1000 appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 50051

ENTRYPOINT ["/app/server"]