-- scripts/post_payload.lua
-- This script configures wrk to send medical telemetry via HTTP POST
-- for performance benchmarking (NFR33).
wrk.method = "POST"
wrk.body   = '{"device_id": "BENCHMARK-001", "heart_rate": 75, "status": "active"}'
wrk.headers["Content-Type"] = "application/json"