FROM golang:1.25-alpine

WORKDIR /app
RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /prepbackend ./cmd

EXPOSE 8080
CMD ["/prepbackend"]