# HelloMix Backend

FROM golang:1.23-alpine AS base

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Development stage with Air for hot reload
FROM base AS development

# Install Air for hot reload and Delve for debugging
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Copy source code
COPY . .

# Expose ports for app and debugger
EXPOSE 8080 2345

# Use Air for hot reload in development
CMD ["air", "-c", ".air.toml"]

# Production builder stage
FROM base AS builder

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
