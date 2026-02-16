# --- Stage 1: Build stage ---
# Using Go 1.25 as requested
FROM golang:1.25-alpine AS builder

# Install ca-certificates in the builder stage
# We will need to copy these to scratch if the app makes HTTPS calls
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Handle dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile a STATIC binary
# CGO_ENABLED=0: Disables C-library linking (crucial for scratch)
# -ldflags="-w -s": Strips debug info and symbol tables to reduce size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/investigation-service/main.go

# --- Stage 2: Final stage (Ultra-minimal) ---
# ... [Previous builder stages] ...

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main

# Define the health check using the binary's own flag
# Interval: check every 30s | Timeout: fail if it takes > 5s
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD ["/main", "-check"]

EXPOSE 8080
EXPOSE 50051

ENTRYPOINT ["/main"]