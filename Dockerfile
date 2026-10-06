# Stage 1: Build con Go 1.26
FROM golang:1.26-alpine AS builder

ENV GOTOOLCHAIN=auto

WORKDIR /app

# Copiar archivos del proyecto
COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ ./pkg/
COPY cmd/ ./cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o orderbook-app ./cmd/api

FROM alpine:3.2

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/orderbook-app .

EXPOSE 8080

ENTRYPOINT ["./orderbook-app"]