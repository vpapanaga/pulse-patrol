// Package app contains the business logic and API handlers for the service.
// It follows the isolation of context principle (ADR #05) defined in the Pulse Patrol architecture.
package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Telemetry represents the data structure received from medical hardware (e.g., EKG sensors).
// This serves as the "Medical Truth" that will be persisted in the Amazon Aurora database.
type Telemetry struct {
	DeviceID  string `json:"device_id"`  // Unique identifier for the hardware equipment
	HeartRate int    `json:"heart_rate"` // Pulse value in beats per minute (BPM)
	Status    string `json:"status"`     // Sensor state (active, error, calibrating)
}

// TelemetryHandler processes HTTP POST requests containing patient vital signs.
// It implements basic method validation and JSON parsing.
//
// Compliance Level: NFR33 (Low Latency < 200ms).
func TelemetryHandler(w http.ResponseWriter, r *http.Request) {
	// Security: Only allow POST methods to prevent exposing medical data via URL parameters.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var data Telemetry
	// Decoding directly from the request body stream for memory efficiency.
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "JSON parsing error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Clinical processing simulation (In next phases, this triggers the Shadow Audit ADR #02).
	fmt.Printf("[INVESTIGATION] Signal received -> Device: %s | HeartRate: %d BPM\n", data.DeviceID, data.HeartRate)

	// Standardized response for on-premise hardware clients.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data successfully received by Pulse Patrol",
	})
}
