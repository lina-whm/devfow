# DevFlow

> Collaborative software development platform — Kanban, Sprints, Analytics, RBAC.

Architecture: `SPA (Next.js) → REST API (Go/Gin) → PostgreSQL + Redis`

## Tech Stack

**Frontend:** React, Next.js (App Router), TypeScript, TailwindCSS, shadcn/ui, TanStack Query, dnd-kit, Recharts, Framer Motion

**Backend:** Go 1.24+, Gin, PostgreSQL 16, Redis 7, JWT, Docker

**Testing:** Go testing, Testcontainers, Playwright, React Testing Library

**Observability:** Prometheus, Grafana, OpenTelemetry (Jaeger)

## Quick Start

```bash
# Clone and run
docker compose -f infrastructure/docker-compose.yml up -d

# Run migrations
cd backend && go run ./cmd/devflow-migrate up

# Access
# Frontend: http://localhost:3000
# API:      http://localhost:8080/api/v1
# Swagger:  http://localhost:8080/swagger/index.html
```

## Features

- Kanban boards with drag-and-drop & optimistic updates
- Sprint planning with burndown charts
- Role-based access control (Owner, Admin, Member)
- Real-time notifications & mentions
- Full-text search with cursor pagination
- Dark theme
- Audit logs & task history

## Project Structure

```
devflow/
├── backend/          # Go API server + workers
│   ├── cmd/          # Entry points
│   ├── internal/
│   │   ├── domain/   # Business entities (DDD)
│   │   ├── application/  # Use cases
│   │   ├── infrastructure/  # DB, Redis, mailer
│   │   ├── api/      # HTTP handlers, middleware
│   │   └── pkg/      # Shared utilities
│   └── docker/
├── frontend/         # Next.js SPA
├── infrastructure/   # Docker Compose, monitoring
└── docs/             # ADRs, API spec, guides
```

## License

MIT
