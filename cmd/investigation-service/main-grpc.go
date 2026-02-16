// Package main initializes the gRPC server for internal service communication.
package main

import (
	"fmt"
	pb "github.com/vpapanaga/pulse-patrol/api/proto"
	"github.com/vpapanaga/pulse-patrol/internal/app"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Listen on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Create a new gRPC server instance
	s := grpc.NewServer()

	// Register our implementation with the gRPC server

	pb.RegisterInvestigationServiceServer(s, &app.GRPCServer{})

	fmt.Println("🏥 Pulse Patrol - Investigation Service")
	fmt.Println("⚡ gRPC Server active on port :50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
