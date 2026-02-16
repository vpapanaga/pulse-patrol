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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	const (
		address     = "localhost:50051"
		concurrency = 10  // Number of parallel workers (like wrk -c)
		requests    = 100 // Total requests per worker
	)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewInvestigationServiceClient(conn)

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("🚀 Starting gRPC Benchmark: %d workers, %d requests each...\n", concurrency, requests)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < requests; j++ {
				_, err := client.SendAlert(context.Background(), &pb.AlertRequest{
					PatientId:    fmt.Sprintf("P-%d", workerID),
					AlertType:    "BENCHMARK",
					CurrentValue: 100,
				})
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

	fmt.Println("\n--- gRPC Benchmark Results ---")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Total Time:     %v\n", duration)
	fmt.Printf("Avg Latency:    %v\n", duration/time.Duration(totalRequests))
	fmt.Printf("Requests/sec:   %.2f\n", float64(totalRequests)/duration.Seconds())
}
