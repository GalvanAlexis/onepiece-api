# ===== STAGE 1: Builder =====
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies for CGO (needed by some packages)
RUN apk add --no-cache git ca-certificates tzdata

# Copy and download dependencies first (layer cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o bin/onepiece-api \
    cmd/api/main.go

# ===== STAGE 2: Runner =====
FROM alpine:3.20

WORKDIR /app

# Security: run as non-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy only the binary and certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/bin/onepiece-api .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./onepiece-api"]
