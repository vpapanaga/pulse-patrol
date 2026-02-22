// Package main initializes the gRPC server for internal service communication.
package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/vpapanaga/pulse-patrol/api/proto"
	"github.com/vpapanaga/pulse-patrol/internal/app"
	"github.com/vpapanaga/pulse-patrol/internal/config" // Import your config package
	"google.golang.org/grpc"
)

func main() {
	// 1. Initialize configuration from .env or environment variables
	config.LoadConfig()

	// 2. Retrieve the gRPC port from config, default to 50051 if not set
	// This matches the externalization strategy used in the REST service
	port := config.GetEnv("GRPC_PORT", "50051")

	// 3. Establish a TCP listener on the configured port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", port, err)
	}

	// 4. Create a new gRPC server instance
	s := grpc.NewServer()

	// 5. Register the Investigation Service implementation with the gRPC server
	pb.RegisterInvestigationServiceServer(s, &app.GRPCServer{})

	fmt.Println("🏥 Pulse Patrol - Investigation Service")
	fmt.Printf("⚡ gRPC Server active on port :%s\n", port)

	// 6. Start the gRPC server and block until the process is terminated
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
