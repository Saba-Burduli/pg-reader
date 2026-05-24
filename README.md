# PG Reader

Minimal, local-first dark-mode reading app for Paul Graham essays.

## Stack
- Backend: Go (`backend/`)
- Frontend: React + Vite (`frontend/`)
- Storage: local JSON cache (`data/articles.json`)

## How scraping works
1. Backend fetches `https://www.paulgraham.com/articles.html`.
2. Extracts essay links and titles.
3. Fetches each essay with rate limiting (`~400ms/request`) and retries (`3` attempts).
4. Converts HTML to readable plain text.
5. Stores structured articles in local cache JSON.

Scraping happens on first backend startup when cache is empty. After that, API serves local cached data for fast reads.

## Run locally

### 1) Backend
```bash
cd backend
go mod tidy
go run .
```

Backend runs on `http://localhost:8080`.

### 2) Frontend
```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173`.

## API
- `GET /articles` -> article list
- `GET /articles/:id` -> full article

## Dev checks
```bash
cd backend && go test ./... && go build ./...
cd frontend && npm run build
```
