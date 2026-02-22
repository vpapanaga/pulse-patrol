// Package main simulates a client calling the Investigation gRPC service.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/vpapanaga/pulse-patrol/api/proto"
	"github.com/vpapanaga/pulse-patrol/internal/config" // Centralized configuration
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Load externalized configuration
	config.LoadConfig()

	// 2. Retrieve the gRPC port from .env, defaulting to 50051
	port := config.GetEnv("GRPC_PORT", "50051")
	targetAddr := fmt.Sprintf("localhost:%s", port)

	fmt.Printf("🧪 Initializing gRPC Service Test on %s...\n", targetAddr)

	// 3. Establish a connection to the server
	// Using insecure credentials for local development as per Lesson 12
	conn, err := grpc.Dial(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC server at %s: %v", targetAddr, err)
	}
	defer conn.Close()

	client := pb.NewInvestigationServiceClient(conn)

	// 4. Perform the RPC call with a 2-second timeout (Resilience pattern)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Mocking a critical heart rate alert
	res, err := client.SendAlert(ctx, &pb.AlertRequest{
		PatientId:    "PATIENT-99",
		AlertType:    "CRITICAL_HEART_RATE",
		CurrentValue: 145,
	})

	if err != nil {
		log.Fatalf("RPC Call failed: %v. Is the gRPC server running on port %s?", err, port)
	}
	fmt.Printf("✅ gRPC Test Success! Status: %s | ID: %s\n", res.Status, res.TrackingId)
}
