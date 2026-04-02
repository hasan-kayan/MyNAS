# Creation Steps of the Authentication Service

---

## 1. Protocol Definition (gRPC Contract)

In this step, we define the **service contract** using Protocol Buffers (proto3).

The `.proto` file describes:

- Service definitions (gRPC RPC methods)
- Request and response message schemas
- Field numbering (for backward compatibility)
- Serialization format

This file represents the **external communication contract** of the Authentication Service.

---

## 1.1 Proto File Location

Create the proto definition file under:

    /api/auth/auth.proto

This structure ensures:

- Clear separation of transport-layer contracts
- Future versioning support (e.g., `v1`, `v2`)
- Clean project organization aligned with Go modules

---

## 1.2 Code Generation

After defining the proto schema, generate the Go bindings using:

```bash
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  api/auth/auth.proto
```

### Explanation of Flags

- `--go_out=.`  
  Generates Go code for message types.

- `--go-grpc_out=.`  
  Generates Go code for gRPC service interfaces.

- `--go_opt=paths=source_relative`  
  Ensures generated files are placed relative to the source `.proto` file.  
  This is required when using Go modules to avoid nested import path folders.

- `--go-grpc_opt=paths=source_relative`  
  Ensures the gRPC service files follow the same relative path rule.

---

## 1.3 Generated Files

The command generates two Go files inside the same directory as the `.proto` file.

---

### `auth.pb.go`

This file contains:

- Generated Go structs for all message types
- Enum definitions
- Protobuf field metadata and tags
- Serialization and deserialization logic
- Proto reflection support

This file represents the **transport data model layer** of the service.

---

### `auth_grpc.pb.go`

This file contains:

- `AuthenticationServiceServer` interface
- `AuthenticationServiceClient` interface
- Default unimplemented server implementation
- gRPC service registration function

Key generated components:

- Server interface (to be implemented by the Authentication Service)
- Client stub (to be used by other services)
- `RegisterAuthenticationServiceServer(...)`

This file defines the **gRPC service abstraction layer**.

---

## Architectural Note

The `.proto` file defines the **external contract boundary** of the Authentication Service.

It does not contain:

- Business logic
- Token validation logic
- Redis integration
- Kafka integration
- Middleware implementations

It strictly defines the communication contract between distributed services and ensures type safety and backward compatibility across the system.
## 2. Determine Non-Functional Requirements

These requirements are not essential for core operations but necessary for service maintainability and observability.

### Logger

For logging, we use the `go/zap` library, which is the best solution for enterprise-grade applications.

### Graceful Shutdown

Graceful shutdown has been added to the main file. During shutdown, it counts down from 10 seconds.

### Health Endpoint

Note: Ordinary `curl` does not work with gRPC due to HTTP/2. Install `grpcurl`:

```bash
brew install grpcurl
```

### Prometheus and Grafana

For observability, Prometheus metrics are collected. The metrics endpoint is configured in `/internal/metrics/metrics.go` and integrated into `main.go`, creating a new metrics port at `9090`.

### Configuration

Go does not read variables from files directly. It checks host environment variables, so the `godotenv` package is required and used in the config package.

**Important:** The main application entry point is in a subdirectory, so there is no root path at the repository level. To run the server correctly, execute it from the `/auth_service` directory, which contains the `.env` file.

