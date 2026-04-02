# RAOS Authentication Service

Enterprise-grade **Authentication & Authorization Microservice** written in **Golang** using **gRPC**, designed with **Clean Architecture**, **SOLID principles**, and a strict **Security-First** mindset.

---

# 📌 Purpose

The RAOS Authentication Service is responsible for:

- User authentication (login / register)
- JWT access & refresh token generation
- Token validation & introspection
- Token revocation & blacklist management
- Multi-Factor Authentication (MFA)
- Service-to-Service authentication
- Secure communication with external User Service

This service does NOT manage user profile data.  
It strictly handles identity and authentication concerns.

---

# 🧠 Architecture Principles

This project strictly follows:

- Clean Architecture
- Domain-Driven Design (DDD)
- SOLID principles
- 12-Factor App methodology
- Security by design
- Separation of concerns

High-Level Layered Architecture:

```
api/                → gRPC proto contracts
cmd/                → application entrypoint
config/             → configuration loader
internal/
    domain/         → entities, interfaces, core rules
    application/    → use cases
    interfaces/     → gRPC handlers
    infrastructure/ → database, redis, jwt, kafka, logger
pkg/                → reusable shared utilities
```

Dependency direction:

```
interfaces → application → domain
infrastructure → domain
```

Domain layer does NOT depend on any external framework.

---

# 🚀 Features

## Core Authentication
- Register
- Login
- Refresh Token
- Logout
- Token Introspection
- Token Revocation

## MFA Support
- Request MFA
- Verify MFA

## Machine Identity
- Service Login (for internal microservices)

## Security
- JWT Access Tokens
- JWT Refresh Tokens
- Redis-backed token blacklist
- bcrypt / argon2 password hashing
- Graceful shutdown support
- Prometheus metrics
- Structured logging (Zap)

---

# ⚙️ Tech Stack

- Go
- gRPC
- PostgreSQL
- Redis
- Kafka (optional domain events)
- Prometheus
- Zap Logger
- Docker-ready

---

# 📦 Configuration

Configuration is environment-based using `.env` or environment variables.

Example:

```env
APP_ENV=development

GRPC_PORT=50051
METRICS_PORT=9090
SHUTDOWN_TIMEOUT=10s

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=raos_auth
DB_SSLMODE=disable

REDIS_ADDRESS=localhost:6379

JWT_SECRET=super_secret_key
JWT_ACCESS_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=720h
```

Configuration is loaded via a centralized `config` package.

---

# 🛠 Running the Service

## 1️⃣ Install Dependencies

```
go mod tidy
```

## 2️⃣ Generate gRPC Code

```
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  api/auth/auth.proto
```

## 3️⃣ Run Service

```
go run cmd/server/main.go
```

---

# 📊 Observability

## Health Endpoints

- /health/live
- /health/ready

## Metrics

Prometheus metrics exposed at:

```
:9090/metrics
```

Metrics include:

- Request latency
- Error counts
- Authentication success/failure rates
- Token issuance counters

---

# 🔐 Security Guidelines

Production environment must enforce:

- Strong JWT secrets
- Encrypted DB connections (SSL required)
- Secret management via secure vault
- Rate limiting (via API Gateway)
- mTLS for service-to-service communication
- Short-lived access tokens
- Refresh token rotation

Never commit secrets to version control.

---

# 🧪 Testing

Run unit tests:

```
go test ./...
```

Integration tests require:

- PostgreSQL instance
- Redis instance

---

# 🔄 Inter-Service Communication

- This service does NOT manage extended user data.
- Communicates with User Service for profile synchronization.
- Can publish domain events (e.g., UserRegistered) to Kafka.

---

# 📈 Production Readiness Checklist

- [ ] Structured logging (Zap)
- [ ] Graceful shutdown
- [ ] Prometheus metrics
- [ ] Dockerfile
- [ ] CI/CD pipeline
- [ ] Secure secret management
- [ ] Load testing
- [ ] Token rotation policy
- [ ] Redis high availability
- [ ] Database migrations strategy

---

# 🏗 Future Enhancements

- OAuth2 support
- OpenID Connect compatibility
- JWKS endpoint
- Key rotation mechanism
- Distributed tracing (OpenTelemetry)
- Audit logging
- Rate limiting middleware

---

# 👨‍💻 Maintainer

RAOS Engineering Team  
Enterprise-grade secure microservices.

---

# 📜 License

Private – RAOS Internal Use Only