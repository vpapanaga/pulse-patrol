-- scripts/post_provision.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

-- We don't necessarily need a body because our Go code
-- generates data based on timestamps, but for a valid REST call:
wrk.body   = "{}"