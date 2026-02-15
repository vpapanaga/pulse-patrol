// Package main in the scripts directory simulates a hardware client (Mock Equipment).
// It is used for local integration testing and data flow validation.
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// Local endpoint for the investigation service ingestion path.
	const apiURL = "http://localhost:8080/v1/telemetry"

	// Mock payload representing data sent by a real-world medical sensor.
	payload := []byte(`{"device_id": "ECG-PATROL-01", "heart_rate": 82, "status": "active"}`)

	fmt.Println("🧪 Initializing Hardware Ingestion Test (HTTP)...")

	// Run a loop to verify server stability and sequential request handling.
	for i := 1; i <= 3; i++ {
		// Sending data via POST request.
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Printf("Test %d failed: %v\n", i, err)
			continue
		}

		// Read the response body to confirm processing by the backend.
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Test %d: Success! Server responded: %s\n", i, string(body))

		// Best Practice: Always close the response body to prevent resource leaks.
		resp.Body.Close()

		time.Sleep(1 * time.Second)
	}
}
