package main

import (
	"fmt"
	"log"
	"net/http"
	"pulse-patrol/internal/app"
)

func main() {
	// Register the handler defined in internal/app
	http.HandleFunc("/v1/telemetry", app.TelemetryHandler)

	fmt.Println("🏥 Pulse Patrol - Investigation Service")
	fmt.Println("🚀 REST Server listening on http://localhost:8080")

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
