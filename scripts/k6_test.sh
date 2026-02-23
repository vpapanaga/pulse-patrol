#!/bin/bash

# Navigate to the project root relative to the script location
cd "$(dirname "$0")/.."

# Parse REST_PORT from .env with macOS/Linux compatibility
PORT=$(grep '^REST_PORT=' .env | cut -d '=' -f2 | tr -d '\r' | xargs)
PORT=${PORT:-8090}

# Verify if the test file exists before starting
if [ ! -f "./scripts/load_test.js" ]; then
    echo "❌ Error: Could not find ./scripts/load_test.js"
    echo "Current working directory: $(pwd)"
    exit 1
fi

echo "🚀 Starting k6 Load Test"
echo "📍 Target: http://localhost:$PORT/v1/provision"
echo "👥 Configuration: 50 Concurrent Virtual Users (Ramp-up)"

# Execute k6 and inject the port as an environment variable
k6 run -e REST_PORT=$PORT ./scripts/load_test.js