
# Stage 1: Build the Go binary

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/ecommerce-server .

# Stage 2: Minimal runtime image

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the compiled binary from builder
COPY --from=builder /app/ecommerce-server .

# Copy .env file (will be overridden by docker-compose env if provided)
COPY .env .env

EXPOSE 8000

ENTRYPOINT ["./ecommerce-server"]
