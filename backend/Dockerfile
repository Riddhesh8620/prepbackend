# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Adjust path if needed (cmd/prepbackend)
RUN go build -o prepbackend ./cmd

# ---------- Runtime stage ----------
FROM alpine:3.20

WORKDIR /app

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/prepbackend /prepbackend

EXPOSE 8080

ENTRYPOINT ["/prepbackend"]