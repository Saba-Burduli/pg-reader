# PG Reader Agent Notes

## Scope
- Build a local-first Paul Graham reader with Go backend and React frontend.
- Keep implementation minimal, fast, and production-clean without overengineering.

## Tech + Conventions
- Backend: Go, standard library HTTP server, JSON file storage under `data/`.
- Frontend: React + Vite, dark-mode-only UI.
- No authentication, no cloud dependencies, no runtime scraping in request path.

## Commands
- Backend run: `cd backend && go run .`
- Frontend run: `cd frontend && npm install && npm run dev`
- Scrape/sync: server auto-syncs if cache missing; manual sync endpoint is not exposed.

## Code Organization
- `backend/models`: shared data types.
- `backend/services`: scraping, storage, and business logic.
- `backend/handlers`: HTTP handlers.
- `backend/main.go`: wiring and server bootstrap.

## Quality Gates
- Always run `go test ./...` and `go build ./...` in `backend/`.
- Always run `npm run build` in `frontend/`.
- Keep README current when behavior or setup changes.
