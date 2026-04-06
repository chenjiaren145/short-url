.PHONY: run run-memory run-redis test build clean compose-up compose-down

# Backend
run: run-memory

run-memory:
	go -C backend run cmd/short-url/main.go -store memory

run-redis:
	docker compose -f docker-compose.yml -f docker-compose.debug.yml up -d redis
	go -C backend run cmd/short-url/main.go -store redis -redis-addr localhost:6379

test:
	go -C backend test ./...

build:
	go -C backend build -o ../short-url cmd/short-url/main.go

# Frontend
frontend-install:
	cd frontend && pnpm install

frontend-dev:
	cd frontend && pnpm dev

# Docker
compose-up:
	docker compose up --build

compose-down:
	docker compose down

clean:
	rm -f short-url
