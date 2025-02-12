FROM golang:1.23.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /api cmd/api/main.go

FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=builder /api /api

CMD ["/api"]