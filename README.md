# PG Reader

Minimal, dark-mode reading library for Paul Graham essays with local progress tracking.

GitHub: https://github.com/Saba-Burduli/pg-reader

## Features
- Scrapes and caches essays from `https://www.paulgraham.com/articles.html`
- Stores title, source URL, publication date (extracted or inferred), and content
- Tracks `read/unread` state per essay (persisted in local JSON cache)
- Date-sorted library with `Unread` and `Completed` sections
- Reading progress dashboard:
  - total essays
  - read essays
  - remaining essays
  - completion percentage bar
- Auto-marks an essay as read when reader scroll reaches the end

## Stack
- Backend: Go (`backend/`)
- Frontend: React + Vite (`frontend/`)
- Storage: local JSON file (`data/articles.json`)

## Scraping Process
1. Fetches the essays index page.
2. Extracts essay links and titles.
3. Fetches each essay with:
   - rate limiting (`~400ms/request`)
   - retry logic (`3` attempts)
4. Converts HTML to readable plain text.
5. Extracts publication date from article body when possible.
6. Falls back to deterministic inferred date when explicit date is missing.
7. Writes structured cache to `data/articles.json`.

Backend uses cache-first behavior and only re-syncs when cache is missing or lacks metadata fields.

## Local Setup

### Backend
```bash
cd backend
go mod tidy
go run .
```
Runs on `http://localhost:8080`.

### Frontend
```bash
cd frontend
npm install
npm run dev
```
Runs on `http://127.0.0.1:5173`.

## API
- `GET /articles`
- `GET /articles/:id`
- `PATCH /articles/:id/read` with body:
```json
{ "isRead": true }
```

## Progress Tracking
Read state is stored in backend cache JSON (`IsRead` field per essay).  
This keeps progress lightweight and local-first without authentication.

## Deployment (Free)
- Public app: https://frontend-eight-delta-64.vercel.app
- Host: Vercel Free tier
- Mode: static data (`frontend/public/articles.json`) + browser localStorage for read progress

This deployment is fully free and does not require a paid backend service.

## Verification Commands
```bash
cd backend && go test ./... && go build ./...
cd frontend && npm run build
```
