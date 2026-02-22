# STEP 1: Build Stage
FROM golang:1.25.5-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
# CGO_ENABLED=0 ensures a static binary for a smaller, more portable image
RUN CGO_ENABLED=0 GOOS=linux go build -o investigation-service ./cmd/investigation-service/main.go

# STEP 2: Final Run Stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates in case the service makes HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/investigation-service .

# Copy the .env file (Optional: usually injected via Docker Compose or K8s)
# COPY .env .

# EXPOSE is documented but overridden by the REST_PORT environment variable
EXPOSE 8080

# HEALTHCHECK instruction uses the internal flag we just implemented
# --interval: how often to check
# --timeout: how long to wait for a response
# --retries: number of failures before marking as UNHEALTHY
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD ./investigation-service --check

# Start the application
CMD ["./investigation-service"]