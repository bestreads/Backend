# Build stage
FROM golang:1.26-rc-alpine3.23 AS builder

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

# ---

# Production stage
FROM busybox:1.37.0-musl

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /out/godocker /usr/local/bin/bestreads

# Copy certificates for HTTPS/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


RUN addgroup -S appuser
RUN adduser -G appuser -D -H -S -s /sbin/nologin appuser

RUN mkdir -p /app/store
RUN chown -R appuser:appuser /app

USER appuser


EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/bestreads" ]
