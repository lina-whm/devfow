# Deployment Guide

## Prerequisites

- Docker & Docker Compose v2
- Go 1.24+ (for local development)
- Node.js 20+ (for frontend development)

## Docker Compose (Development)

```bash
# Start all services
docker compose -f infrastructure/docker-compose.yml up -d

# Run database migrations
cd backend && go run ./cmd/devflow-migrate up

# Access the application
open http://localhost:3000
```

## Manual Deployment

### Backend

```bash
cd backend
go build -o bin/devflow ./cmd/devflow
DATABASE_URL=postgres://user:pass@host:5432/devflow REDIS_URL=redis://host:6379 ./bin/devflow
```

### Frontend

```bash
cd frontend
npm run build
npm start
```

## Production Considerations

### SSL/TLS
- Terminate SSL at the reverse proxy (Nginx/Caddy/Traefik)
- Use Let's Encrypt for automatic certificate renewal
- Set `SameSite=Strict; Secure` on cookies

### Database
- Use managed PostgreSQL (Aiven, Supabase, AWS RDS)
- Enable automated backups (pg_dump or WAL archiving)
- Set up connection pooling (PgBouncer) if needed

### Redis
- Use Redis 7 with persistence (AOF + RDB)
- Set maxmemory-policy to allkeys-lru for cache

### Security
- Rotate JWT secrets regularly
- Use environment variables or a secrets manager (HashiCorp Vault, Doppler)
- Enable rate limiting at the reverse proxy level
- Run the application as a non-root user (provided in Dockerfile)

### Monitoring
- Prometheus + Grafana for metrics (docker-compose.monitoring.yml)
- Jaeger for distributed tracing
- Structured JSON logs shipped to your logging platform

### Scaling
- Horizontally scale the backend behind a load balancer
- Use Redis as a shared session store (already implemented)
- Ensure the database connection pool is sized appropriately
