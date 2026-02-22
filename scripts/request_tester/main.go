// Package main in the scripts directory simulates a hardware client (Mock Equipment).
// It is used for local integration testing and data flow validation.
package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2" // Modern random generation for Go 1.22+
	"net/http"
	"time"

	"github.com/vpapanaga/pulse-patrol/internal/config"
)

func main() {
	// 1. Initialize configuration from .env or environment variables
	config.LoadConfig()

	// 2. Dynamically build the endpoint using the REST_PORT from .env
	port := config.GetEnv("REST_PORT", "8080")
	apiURL := fmt.Sprintf("http://localhost:%s/v1/telemetry", port)

	fmt.Printf("🧪 Initializing Hardware Ingestion Test (HTTP) on port %s...\n", port)

	// Run a loop to verify server stability and sequential request handling.
	for i := 1; i <= 5; i++ {
		// 3. Generate random medical telemetry data
		// Simulating a heart rate between 60 and 120 BPM
		heartRate := 60 + rand.IntN(61)

		// Create a dynamic payload for each iteration
		payload := []byte(fmt.Sprintf(
			`{"device_id": "ECG-PATROL-01", "heart_rate": %d, "status": "active"}`,
			heartRate,
		))

		// Sending data via POST request
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Printf("Test %d failed: %v (Is the server running?)\n", i, err)
			continue
		}

		// Read the response body to confirm processing by the backend
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Test %d: [HR: %d] Success! Server responded: %s\n", i, heartRate, string(body))

		// Best Practice: Always close the response body to prevent resource leaks
		resp.Body.Close()

		// Simulate a sampling interval (medical sensors usually sample at fixed rates)
		time.Sleep(1 * time.Second)
	}
}
