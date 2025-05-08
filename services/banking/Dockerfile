FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum* ./

# Copy the go-core package
COPY ../packages/go-core /packages/go-core

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /banking-service ./cmd/server

# Use a minimal alpine image for the final image
FROM alpine:3.18

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /banking-service /banking-service

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Run as non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose the service port
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/banking-service"]