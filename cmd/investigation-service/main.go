package main

import (
	"flag"
	"fmt"
	"github.com/vpapanaga/pulse-patrol/internal/app"
	"net/http"
	"os"
)

func main() {
	// 1. Define a flag for health check mode
	isCheck := flag.Bool("check", false, "Run in health check mode")
	flag.Parse()

	// 2. If in check mode, perform the ping and exit
	if *isCheck {
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1) // Docker interprets non-zero as UNHEALTHY
		}
		os.Exit(0) // Docker interprets zero as HEALTHY
	}

	// 3. Normal Server Startup
	http.HandleFunc("/v1/telemetry", app.TelemetryHandler)

	// Add the health endpoint for the orchestrator/Docker to hit
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("🏥 Pulse Patrol - Investigation Service Active")
	http.ListenAndServe(":8080", nil)
}
