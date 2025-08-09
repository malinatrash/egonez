# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /egonez ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /egonez .

# Copy configuration files
COPY .env.example .env

# Expose port (if needed)
# EXPOSE 8080

# Command to run the executable
CMD ["./egonez"]
