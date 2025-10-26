# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Native build with optimization flags (no GOOS/GOARCH)
RUN CGO_ENABLED=0 go build \
    -ldflags='-w -s' \
    -a -installsuffix cgo \
    -o main .

# Minimal scratch image
FROM scratch

# Import CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/main /

# Expose port
EXPOSE 8000

# Run binary
CMD ["/main"]