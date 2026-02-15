// Package main is the entry point for the investigation microservice.
// It configures the servers and routes required for the service to function within the AWS ecosystem.
package main

import (
	"fmt"
	"github.com/vpapanaga/pulse-patrol/internal/app"
	"log"
	"net/http"
)

func main() {
	// Register the primary route for telemetry ingestion.
	// The handler is imported from internal/app to keep the private business logic protected.
	http.HandleFunc("/v1/telemetry", app.TelemetryHandler)

	fmt.Println("🏥 Pulse Patrol - Investigation Service")
	fmt.Println("📡 HTTP Ingestion Server active on port :8080")
	fmt.Println("---------------------------------------------------")

	// Start the blocking HTTP server. In production (AWS Fargate),
	// this server is typically monitored by an Application Load Balancer (ALB).
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Critical error starting the server: %v", err)
	}
}
