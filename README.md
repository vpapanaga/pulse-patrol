# Pulse-Patrol - Investigation Service

## 🏥 Project Overview
The **Investigation Service** is a high-performance backend component developed as part of the SkilLab Software Architecture course. It acts as the bridge between on-premise medical hardware and the AWS Cloud-native ecosystem, ensuring efficient telemetry ingestion and clinical alert management.

## 🏗 Project Structure
This repository follows the **Standard Go Project Layout** (Option 1) to maintain a clean separation between the transport layer, business logic, and testing utilities.

```text
pulse-patrol/
├── api/
│   └── proto/
│       └── investigation.proto     # gRPC Protocol Buffer definitions
├── cmd/
│   └── investigation-service/
│       └── main.go                 # Main entry point for the service
├── configs/                        # Configuration files (YAML/JSON)
├── deployment/
│   └── Dockerfile                  # Container definition for AWS Fargate
├── internal/
│   ├── app/
│   │   ├── handler.go              # API request handlers
│   │   └── handlers_test.go        # Unit tests for handlers
│   ├── domain/                     # Domain models and business logic
│   └── repository/                 # Data access layer (PostgreSQL/Aurora)
├── pkg/                            # Public, reusable library code
├── scripts/
│   └── request_tester/
│       └── main.go                 # Standalone hardware mock simulator
├── test/
│   └── integration-test.go         # End-to-end integration tests
├── go.mod                          # Go module dependencies
└── README.md                       # Project documentation
```  
## 🏗 Howto run the project
### Start the HTTP Web Server on port 8080
```bash
go run cmd/investigation-service/main.go
``` 
### Start the HTTP Web Clients 
```bash
go run scripts/request_tester/main.go
``` 

### Execute WRK tests
```bash
wrk -t2 -c100 -d30s -s scripts/post_payload.lua http://localhost:8080/v1/telemetry
````
