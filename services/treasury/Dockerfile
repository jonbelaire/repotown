FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o treasury-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /app/treasury-service .
COPY --from=builder /app/migrations ./migrations

# Expose the service port
EXPOSE 8083

# Run the application
CMD ["./treasury-service"]