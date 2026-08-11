# Full-Stack Calculator

A full-stack calculator built with a Go REST API and a React + TypeScript frontend.

## Tech Stack

- Backend: Go 1.26 (`net/http` + `chi` router)
- Frontend: React 19 + TypeScript + Vite
- Frontend testing: Vitest + React Testing Library
- Backend testing: `go test`
- Containerization: Docker + Docker Compose with nginx reverse proxy

## Project Structure

```text
fullstack-calculator/
  backend/
    cmd/server/              # application entrypoint
    internal/calculator/     # pure arithmetic logic
    internal/httpapi/        # routes, handlers, validation, responses
    internal/config/         # env-based runtime config
  frontend/
    src/components/          # UI components
    src/api/                 # API client and types
    src/lib/                 # frontend validation and payload helpers
  docker-compose.yml
```



## Prerequisites

- Go 1.26+
- Node.js 24+
- npm
- Docker Desktop (optional, for Compose)



## Environment Variables



### Backend (`backend/.env`)

Copy `backend/.env.example` to `backend/.env` and adjust if needed.

- `PORT` (default: `8080`)
- `ALLOWED_ORIGINS` (comma-separated origins)



### Frontend (`frontend/.env`)

Copy `frontend/.env.example` to `frontend/.env`.

- Local dev should use: `VITE_API_BASE_URL=http://localhost:8080`
- Docker Compose uses same-origin proxy (`VITE_API_BASE_URL` is injected as empty at build time)

> Never commit real `.env` files.



## Run Locally (without Docker)



### 1) Start backend

```bash
cd backend
go run ./cmd/server
```

Backend runs on `http://localhost:8080`.

### 2) Start frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173`.

## Run with Docker Compose

From repository root:

```bash
docker compose up --build
```

- Frontend: `http://localhost` (port 80)
- Backend: `http://localhost:8080`

The frontend container serves static files with nginx and proxies `/api/*` to the backend service.

Stop stack:

```bash
docker compose down
```



## Test Commands



### Backend

```bash
cd backend
go test -cover ./...
```



### Frontend

```bash
cd frontend
npm test
npm run test:coverage
```



## Coverage Reports



### Backend coverage report (PowerShell)

```powershell
cd backend
go test "-coverprofile=coverage.out" ./...
go tool cover "-html=coverage.out"
```



### Frontend coverage report

```bash
cd frontend
npm run test:coverage
```

Vitest writes coverage output under `frontend/coverage/`.

## API Specification

Base path: `/api/v1`

### Success response format

```json
{ "result": 123.45 }
```



### Error response format

```json
{ "error": "clear description of the problem" }
```



### Endpoints

- `POST /add` with `{ "a": number, "b": number }`
- `POST /subtract` with `{ "a": number, "b": number }`
- `POST /multiply` with `{ "a": number, "b": number }`
- `POST /divide` with `{ "a": number, "b": number }`
- `POST /exponent` with `{ "base": number, "exp": number }`
- `POST /sqrt` with `{ "a": number }`
- `POST /percentage` with `{ "a": number, "b": number }` (a% of b)



## cURL Examples

```bash
# Add
curl -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" \
  -d '{"a":2,"b":3}'

# Subtract
curl -X POST http://localhost:8080/api/v1/subtract \
  -H "Content-Type: application/json" \
  -d '{"a":7,"b":4}'

# Multiply
curl -X POST http://localhost:8080/api/v1/multiply \
  -H "Content-Type: application/json" \
  -d '{"a":3,"b":5}'

# Divide
curl -X POST http://localhost:8080/api/v1/divide \
  -H "Content-Type: application/json" \
  -d '{"a":10,"b":2}'

# Exponent
curl -X POST http://localhost:8080/api/v1/exponent \
  -H "Content-Type: application/json" \
  -d '{"base":2,"exp":8}'

# Square root
curl -X POST http://localhost:8080/api/v1/sqrt \
  -H "Content-Type: application/json" \
  -d '{"a":9}'

# Percentage (a% of b)
curl -X POST http://localhost:8080/api/v1/percentage \
  -H "Content-Type: application/json" \
  -d '{"a":10,"b":200}'
```



## Key Design Decisions

- Thin HTTP handlers: decode/validate -> call pure logic -> write response
- Strict JSON validation (`DisallowUnknownFields`) to avoid silent payload mistakes
- Operation-specific request bodies (e.g. `sqrt` sends only `{a}`)
- Shared numeric bound (`range -1e15 to 1e15`) aligned with float precision constraints
- Frontend input sanitization: numeric-only typing, 15 integer digits, 10 decimals
- nginx proxy in Compose to keep frontend API calls same-origin



## Known Limitations

- No persistent history of operations
- No authentication/authorization (out of scope)



## Next Improvements

- Add E2E browser tests (Playwright/Cypress)
- Add CI workflow for lint + tests + build
- Improve accessibility and keyboard shortcuts in calculator UI

