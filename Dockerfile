# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pricewatcher ./cmd/pricewatcher

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/pricewatcher .
COPY --from=builder /app/migrations ./migrations

# Copy configuration
COPY config.yaml /app/config.yaml

# Expose the application port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./pricewatcher"]
