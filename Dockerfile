# Build stage
FROM golang:1.25 AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server ./cli/mcp

# Runtime stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /build/server /server

# Set working directory
WORKDIR /data

# Expose HTTP port (only used when transport=http)
EXPOSE 8080

# Default command runs in stdio mode
ENTRYPOINT ["/server"]
CMD ["serve"]
