# syntax=docker/dockerfile:1.7

# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/egonez ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/egonez .

# Make the binary executable
RUN chmod +x egonez

# Set the entrypoint
ENTRYPOINT ["/app/egonez"]