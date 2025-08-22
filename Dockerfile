# Multi-stage build for BeesInTheTrap Go application
# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod ./
# Copy go.sum if it exists (conditional copy for projects without external dependencies)
COPY go.su[m] ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o beesinthetrap ./cmd/beesinthetrap

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install ca-certificates for any HTTPS calls (optional but good practice)
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/beesinthetrap .

# Change ownership to non-root user
RUN chown appuser:appgroup beesinthetrap

# Switch to non-root user
USER appuser

# Expose port (not needed for this CLI app, but good practice)
# EXPOSE 8080

# Set the binary as entrypoint
ENTRYPOINT ["./beesinthetrap"]

# Default command (can be overridden)
CMD []
