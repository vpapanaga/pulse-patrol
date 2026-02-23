# STEP 1: Build Stage
FROM golang:1.25.5-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the provisioning binary
# Using CGO_ENABLED=0 for a purely static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o provisioning-service ./cmd/provisioning/main.go

# STEP 2: Runtime Stage
FROM alpine:latest

# Security: HIPAA/NFR36 compliance using non-root user
RUN adduser -D pulseuser
WORKDIR /home/pulseuser/

# Install CA certificates for secure DB connections
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder
COPY --from=builder /app/provisioning-service .

# Port Documentation
# The service will use the REST_PORT defined in your .env
EXPOSE 8081

USER pulseuser

# HEALTHCHECK: Verify service availability on the assigned port
# Uses wget (built-in Alpine) to check the /health endpoint
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Launch the service
ENTRYPOINT ["./provisioning-service"]