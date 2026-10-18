# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/montra-api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S montra && adduser -S -G montra montra

WORKDIR /app
COPY --from=builder /out/montra-api /app/montra-api
COPY --from=builder /src/migrations /app/migrations

USER montra
EXPOSE 8080

ENTRYPOINT ["/app/montra-api"]
