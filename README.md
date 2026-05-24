# PG Reader

A minimal dark-mode reading app for Paul Graham essays.

## Live App
https://frontend-eight-delta-64.vercel.app

## What It Does
- Loads essays from Paul Graham’s official essays page
- Sorts essays by extracted publication date
- Shows the official publication date when available
- Lets readers mark essays as read/unread
- Persists progress locally for returning sessions

## Tech Stack
- Backend: Go
- Frontend: React + Vite
- Storage:
  - Local JSON cache for scraped essay data
  - Local browser storage for read progress (in deployed static mode)

## Local Development

### 1) Backend
```bash
cd backend
go mod tidy
go run .
```

### 2) Frontend
```bash
cd frontend
npm install
npm run dev
```

Frontend: `http://127.0.0.1:5173`  
Backend: `http://localhost:8080`

## API (Backend Mode)
- `GET /articles`
- `GET /articles/:id`
- `PATCH /articles/:id/read`

## Data Source
- https://www.paulgraham.com/articles.html

## License
MIT
