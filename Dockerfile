# Build stage
FROM golang:1-alpine AS builder

WORKDIR /app

# Install git and SSL certificates
RUN apk add --no-cache git ca-certificates

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the project
COPY . .

# Explicitly create out dir
RUN mkdir -p /out

# Build Go API
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags="-s -w" -o /out/godocker ./

# Create non-root user
RUN adduser -D -H -h /nonexistent -s /sbin/nologin -u 10001 appuser && \
    grep appuser /etc/passwd > /out/passwd

# ---

# Production stage
FROM scratch

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /out/godocker /usr/local/bin/godocker

# Copy certificates for HTTPS/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy non-root user
COPY --from=builder /out/passwd /etc/passwd
USER appuser

RUN mkdir -p /app/store
RUN chown appuser /app/store

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/godocker" ]
