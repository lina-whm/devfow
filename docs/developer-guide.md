# Developer Guide

## Setup

```bash
# Clone the repository
git clone https://github.com/your-org/devflow.git
cd devflow

# Start infrastructure (PostgreSQL, Redis, MinIO)
docker compose -f infrastructure/docker-compose.yml up -d postgres redis minio

# Install backend dependencies
cd backend
go mod download

# Run migrations
go run ./cmd/devflow-migrate up

# Start backend (hot reload requires air)
go run ./cmd/devflow

# In another terminal, start frontend
cd frontend
npm install
npm run dev
```

## Environment Variables

Copy `.env.example` to `.env` and adjust as needed. Key variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://devflow:devflow@localhost:5432/devflow?sslmode=disable` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379/0` |
| `JWT_SECRET` | HMAC-SHA256 key for access tokens | (required) |
| `JWT_REFRESH_SECRET` | HMAC-SHA256 key for refresh tokens | (required) |
| `CORS_ORIGINS` | Allowed CORS origins | `http://localhost:3000` |
| `LOG_LEVEL` | Log level: debug/info/warn/error | `info` |

## Available Commands

```bash
make backend        # Start Go server
make frontend       # Start Next.js dev server
make dev            # Start all services with Docker Compose
make test-backend   # Run Go tests
make test-frontend  # Run frontend tests
make lint-backend   # Lint Go code
make migrate        # Run database migrations
make swagger        # Generate Swagger docs
```

## Architecture Overview

```
┌──────────────┐     ┌──────────────────┐     ┌────────────┐
│   Browser    │────▶│  Go/Gin Server   │────▶│ PostgreSQL │
│  (Next.js)   │     │  Port 8080       │     │  Port 5432 │
└──────────────┘     └────────┬─────────┘     └────────────┘
                              │                      │
                              ▼                      ▼
                        ┌────────────┐     ┌────────────────┐
                        │   Redis    │     │   Background   │
                        │  Port 6379 │     │    Workers     │
                        └────────────┘     └────────────────┘
```

## Testing

```bash
# Run all tests
cd backend && go test ./... -race -count=1

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Integration tests (require Docker)
go test ./test/... -tags=integration
```
