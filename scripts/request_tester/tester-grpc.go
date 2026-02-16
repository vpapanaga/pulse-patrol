// Package main simulates a client calling the Investigation gRPC service.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/vpapanaga/pulse-patrol/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Establish a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewInvestigationServiceClient(conn)

	fmt.Println("🧪 Initializing gRPC Service Test...")

	// Perform the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SendAlert(ctx, &pb.AlertRequest{
		PatientId:    "PATIENT-99",
		AlertType:    "CRITICAL_HEART_RATE",
		CurrentValue: 145,
	})

	if err != nil {
		log.Fatalf("RPC Call failed: %v", err)
	}

	fmt.Printf("✅ gRPC Test Success! Status: %s | ID: %s\n", res.Status, res.TrackingId)
}
