// Package main provides a high-concurrency benchmark for the gRPC Investigation Service.
// It simulates multiple concurrent clinical orchestrators to test NFR33 compliance.
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/vpapanaga/pulse-patrol/api/proto"
	"github.com/vpapanaga/pulse-patrol/internal/config" // Import centralized config
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Load configuration from .env or environment variables
	config.LoadConfig()

	// 2. Retrieve gRPC address dynamically
	port := config.GetEnv("GRPC_PORT", "50051")
	address := fmt.Sprintf("localhost:%s", port)

	const (
		concurrency = 10  // Number of parallel workers (Simulating multiple devices)
		requests    = 100 // Total requests per worker
	)

	// 3. Establish the gRPC connection
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Benchmark failed - could not connect to %s: %v", address, err)
	}
	defer conn.Close()
	client := pb.NewInvestigationServiceClient(conn)

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("🚀 Starting gRPC Benchmark on %s: %d workers, %d requests each...\n", address, concurrency, requests)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < requests; j++ {
				// We use a context with a timeout for each call to prevent hanging workers
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := client.SendAlert(ctx, &pb.AlertRequest{
					PatientId:    fmt.Sprintf("P-%d", workerID),
					AlertType:    "BENCHMARK",
					CurrentValue: 100,
				})
				cancel() // Release context resources

				if err != nil {
					fmt.Printf("Worker %d failed: %v\n", workerID, err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	totalRequests := concurrency * requests
	// 4. Output results in a format useful for ARD NFR33 compliance reporting
	fmt.Println("\n--- gRPC Benchmark Results ---")
	fmt.Printf("Target Service: %s\n", address)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Total Time:     %v\n", duration)
	fmt.Printf("Avg Latency:    %v\n", duration/time.Duration(totalRequests))
	fmt.Printf("Requests/sec:   %.2f\n", float64(totalRequests)/duration.Seconds())
}
