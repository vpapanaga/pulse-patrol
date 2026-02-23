// Package main initializes the Provisioning Service for Pulse Patrol.
// This service is specifically designed for high-performance data seeding
// and benchmarking against the Aurora PostgreSQL schema.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vpapanaga/pulse-patrol/internal/config"
)

// ProvisioningServer manages the database connection pool and request handling.
type ProvisioningServer struct {
	db *pgxpool.Pool
}

func main() {
	// 1. Initialize centralized configuration from .env or system environment
	config.LoadConfig()

	// 2. Setup the PostgreSQL/Aurora connection pool
	// We use pgxpool for concurrency management and automatic connection recycling.
	dsn := config.GetDatabaseDSN()
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Critical Error: Unable to parse DB DSN: %v", err)
	}

	// Performance Tuning: Optimizing for 2026 Medical-Grade scalability (NFR33)
	dbConfig.MaxConns = 30                     // Matches expected peak worker threads
	dbConfig.MinConns = 5                      // Keep baseline connections active
	dbConfig.MaxConnIdleTime = 1 * time.Minute // Recycle inactive connections

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Critical Error: Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// 3. Initialize the server controller
	server := &ProvisioningServer{db: pool}

	// 4. Register HTTP Routes
	mux := http.NewServeMux()

	// Primary endpoint for benchmarking (wrk, k6)
	mux.HandleFunc("/v1/provision", server.handleProvisioning)

	// Readiness/Liveness probe for SigNoz and Docker orchestration
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HEALTHY"))
	})

	// 5. Port Configuration Logic
	// Defaulting to 8081 for provisioning to avoid port collision with the main app
	serverAddr := ":" + config.GlobalConfig.RestPort
	if serverAddr == ":8080" || serverAddr == ":" {
		serverAddr = ":8081"
	}

	fmt.Printf("🏥 Pulse Patrol - Provisioning Service Online\n")
	fmt.Printf("📊 Benchmarking URL: http://localhost%s/v1/provision\n", serverAddr)
	fmt.Printf("🔗 Database Target: %s\n", config.GlobalConfig.DBHost)

	// Start the HTTP Server
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleProvisioning manages the creation of institutions and devices within a single transaction.
// It ensures Atomicity: either both records are created, or neither is.
func (s *ProvisioningServer) handleProvisioning(w http.ResponseWriter, r *http.Request) {
	// Enforce POST method for data modification
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use request context for timeout/cancellation propagation to the DB
	ctx := r.Context()

	// Begin a new database transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Printf("Transaction Start Error: %v", err)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	// Ensure transaction is closed via rollback if Commit is never reached
	defer tx.Rollback(ctx)

	// Data Generation for Benchmark Variety
	startTime := time.Now()
	uniqueID := startTime.UnixNano()
	clinicName := fmt.Sprintf("Clinic %d", uniqueID)
	licenseKey := fmt.Sprintf("LIC-%d", rand.IntN(999999))
	serialNumber := fmt.Sprintf("SN-%d", uniqueID)

	// STEP 1: Insert Institution record (NFR36: Tenant Isolation)
	var institutionID string
	err = tx.QueryRow(ctx,
		"INSERT INTO institutions (name, license_key) VALUES ($1, $2) RETURNING institution_id",
		clinicName, licenseKey,
	).Scan(&institutionID)
	if err != nil {
		log.Printf("DB Error (Institution): %v", err)
		http.Error(w, "Failed to create institution resource", http.StatusUnprocessableEntity)
		return
	}

	// STEP 2: Insert Device record (FR76: Device Management)
	_, err = tx.Exec(ctx,
		"INSERT INTO devices (serial_number, device_model, status) VALUES ($1, $2, $3)",
		serialNumber, "Pulse-V1-Sensor", "AVAILABLE",
	)
	if err != nil {
		log.Printf("DB Error (Device): %v", err)
		http.Error(w, "Failed to create device resource", http.StatusUnprocessableEntity)
		return
	}

	// Finalize the database changes
	if err := tx.Commit(ctx); err != nil {
		log.Printf("DB Commit Error: %v", err)
		http.Error(w, "Finalization failed", http.StatusInternalServerError)
		return
	}

	// 6. Response Construction for Observability
	// Including latency_ms allows SigNoz to correlate HTTP logs with DB performance.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "provisioned",
		"institution_id": institutionID,
		"device_serial":  serialNumber,
		"latency_ms":     time.Since(startTime).Milliseconds(),
	})
}
