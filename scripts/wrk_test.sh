#!/bin/bash

# 1. Load the port from the .env file
# This extracts the REST_PORT value and cleans any Windows-style line endings (\r)
PORT=$(grep '^REST_PORT=' .env | cut -d '=' -f2 | sed 's/\r//' | xargs)

# 2. Set fallback to 8090 (matching your .env configuration)
PORT=${PORT:-8090}

# 3. Define the Provisioning endpoint
ENDPOINT="/v1/provision"
HEALTH_ENDPOINT="/health"

echo "🚀 Pulse Patrol - Provisioning Benchmark"
echo "📍 Target: http://localhost:$PORT$ENDPOINT"

# 4. Pre-flight Check: Ensure the service is healthy and the DB is connected
# We check the /health endpoint first to ensure the container is ready
echo "🔍 Checking service status..."
if ! curl -s "http://localhost:$PORT$HEALTH_ENDPOINT" | grep "HEALTHY" > /dev/null; then
    echo "❌ Error: Provisioning service is not HEALTHY on port $PORT."
    echo "💡 Tip: Run 'docker compose logs provisioning-service' to check for DB connection errors."
    exit 1
fi

# 5. Execute wrk with the POST script
# -t12: 12 threads
# -c400: 400 concurrent connections
# -d30s: duration of 30 seconds
# -s: path to the Lua script for POST requests
echo "🔥 Starting high-load test (400 connections)..."

wrk -t12 -c400 -d30s -s scripts/post_provision.lua "http://localhost:$PORT$ENDPOINT"