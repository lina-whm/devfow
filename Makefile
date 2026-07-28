.PHONY: dev up down backend frontend lint test migrate

dev:
	docker compose -f infrastructure/docker-compose.yml -f infrastructure/docker-compose.dev.yml up

up:
	docker compose -f infrastructure/docker-compose.yml up -d

down:
	docker compose -f infrastructure/docker-compose.yml down

backend:
	cd backend && go run ./cmd/devflow

worker:
	cd backend && go run ./cmd/devflow-worker

migrate:
	cd backend && go run ./cmd/devflow-migrate up

migrate-down:
	cd backend && go run ./cmd/devflow-migrate down

frontend:
	cd frontend && npm run dev

lint-backend:
	cd backend && golangci-lint run

lint-frontend:
	cd frontend && npm run lint

test-backend:
	cd backend && go test ./... -race -count=1 -coverprofile=coverage.out

test-frontend:
	cd frontend && npm run test

test-e2e:
	cd frontend && npx playwright test

build-backend:
	cd backend && go build -o bin/devflow ./cmd/devflow

build-frontend:
	cd frontend && npm run build

swagger:
	cd backend && swag init -g cmd/devflow/main.go -o internal/api/docs

coverage:
	cd backend && go tool cover -html=coverage.out -o coverage.html

docker-build:
	docker compose -f infrastructure/docker-compose.yml build
