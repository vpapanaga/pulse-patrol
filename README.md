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
### Start the gRPC Server on port 50051
```bash
go run cmd/investigation-service/main-grpc.go
``` 
### Start the gRPC Client
```bash
go run scripts/request_tester/tester-grpc.go
```
### Execute Bench tests with gRPC
```bash
go run scripts/request_tester/bench_grpc.go
````
### Execute Bench tests with ghz

#### Install the ghz application 
```bash
# MacOS
brew install ghz
# Go Install
go install github.com/bojand/ghz/cmd/ghz@latest
````
#### Run the tests
```bash
ghz --insecure --proto ./api/proto/investigation.proto --call investigation.InvestigationService.SendAlert --data '{"patient_id": "GHZ-001", "alert_type": "LOAD_TEST", "current_value": 95}' -n 1000 -c 50 localhost:50051
````
## Using the dockerized container
### Create the local container
```bash
docker build -t pulse-patrol-investigation -f deployment/Dockerfile .
````
### Execute the container
```bash
docker run -p 8080:8080 -p 50051:50051 pulse-patrol-investigation
````
