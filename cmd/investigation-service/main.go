package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/vpapanaga/pulse-patrol/internal/app"
	"github.com/vpapanaga/pulse-patrol/internal/config" // Ensure this path matches your go.mod
)

func main() {
	// 1. Initialize configuration from .env or environment variables
	config.LoadConfig()

	// 2. Define a flag for health check mode (used by Docker HEALTHCHECK)
	isCheck := flag.Bool("check", false, "Run in health check mode to verify service status")
	flag.Parse()

	// Retrieve the REST port from config, default to 8080 if not set
	port := config.GetEnv("REST_PORT", "8080")
	healthAddr := fmt.Sprintf("http://localhost:%s/health", port)

	// 3. Health Check Mode Execution
	// When Docker runs 'investigation-service --check', this block executes
	if *isCheck {
		resp, err := http.Get(healthAddr)
		if err != nil || resp.StatusCode != http.StatusOK {
			// Non-zero exit tells the orchestrator the service is UNHEALTHY
			os.Exit(1)
		}
		// Zero exit tells the orchestrator the service is HEALTHY
		os.Exit(0)
	}

	// 4. Standard API Routes Setup
	// Primary telemetry ingestion endpoint
	http.HandleFunc("/v1/telemetry", app.TelemetryHandler)

	// Liveness/Readiness probe endpoint for Docker/Kubernetes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 5. Normal Server Startup
	fmt.Printf("🏥 Pulse Patrol - Investigation Service Active on Port %s\n", port)

	// Start the HTTP server. The address is now dynamic based on the config.
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Critical failure starting server: %v\n", err)
		os.Exit(1)
	}
}
