# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo short URL service with:
- **Backend**: Go + Gin API (`backend/`)
- **Frontend**: React + Vite (`frontend/`)

## Development Commands

### Backend
```bash
# Run backend with in-memory storage (fastest for local dev)
make run-memory

# Run backend with Redis storage
make run-redis

# Run tests
make test

# Build binary
make build
```

Backend flags:
- `-store`: "memory" or "redis" (default: memory)
- `-redis-addr`: Redis address (default: localhost:6379)
- `-port`: HTTP server port (default: 8081)
- `-base-url`: Public base URL for short links (e.g., https://sho.rt)

### Frontend
```bash
cd frontend
npm ci              # Install dependencies
npm run dev         # Start dev server (proxies /api to localhost:8081)
npm run build       # Build for production
npm run typecheck   # TypeScript type checking
npm run lint        # ESLint
```

### Docker
```bash
make compose-up                 # Start app + Redis (Redis not exposed)
make compose-up-debug-redis     # Start with Redis exposed on 6379 for debugging
make compose-down               # Stop all services
```

## Architecture

### Backend Structure
```
backend/
├── cmd/short-url/main.go          # Entry point, wire dependencies
├── internal/
│   ├── api/handler.go              # HTTP handlers (Gin)
│   ├── service/service.go          # Business logic layer
│   ├── shortener/generator.go      # Short code generation
│   └── store/
│       ├── store.go                # Store interface
│       ├── memory.go               # In-memory implementation
│       └── redis.go                # Redis implementation
```

**Layered architecture**:
1. **API layer** (`api/handler.go`): HTTP request/response handling with Gin
2. **Service layer** (`service/service.go`): Business logic, uses Store interface
3. **Store layer** (`store/`): Storage abstraction with interface - memory or Redis

The Store interface (`Save`, `Load`, `IncrementVisits`, `GetVisits`) allows swapping storage backends without changing service logic. Service layer handles short code generation and coordinates store operations.

### Frontend Structure
```
frontend/src/
├── components/
│   ├── URLForm.tsx      # Form to create short URLs
│   └── URLList.tsx      # Display created URLs with stats
├── services/api.ts      # API client (posts to /shorten, fetches /stats)
├── types/url.ts         # TypeScript types
└── main.tsx             # Entry point
```

Vite dev server proxies `/api` requests to `localhost:8081` for backend calls.

## API Endpoints

- `GET /` - Welcome message
- `GET /healthz` - Health check
- `GET /readyz` - Readiness check
- `POST /shorten` - Create short URL (JSON: `{"original_url": "https://..."}`)
- `GET /:shortCode` - Redirect to original URL
- `HEAD /:shortCode` - Redirect (HEAD method)
- `GET /:shortCode/stats` - Get visit count

## Testing

Backend tests follow Go conventions: `*_test.go` files next to source files. Run with `make test` or `go test ./...` from the backend directory.
