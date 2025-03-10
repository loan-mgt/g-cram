# Build stage
FROM golang:1.24.1-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY cramer/go.mod cramer/go.sum* ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY cramer/ ./

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o cramer ./cmd/main.go

# Final stage
FROM alpine:3.21.3

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata ffmpeg

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/cramer .

# Expose API port
EXPOSE 8080

# Use non-root user
USER appuser

# Run the application
CMD ["./cramer"]