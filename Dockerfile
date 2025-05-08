# Etapa de build con Alpine (ligera)
FROM golang:1.24.3-alpine3.21 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilación estática para Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o products ./cmd/main.go

# Imagen final ultra mínima
FROM alpine:3.21

WORKDIR /app
COPY --from=builder /app/products .

CMD ["./products"]
