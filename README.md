# DevFlow OS v1

A personal operating system for dev creators: daily work → content → opportunities → money.

## Stack

- **Frontend**: React (Vite) + Tailwind CSS + shadcn/ui — pastel theme
- **Backend**: Go (clean architecture) + PostgreSQL
- **AI**: Gemini API (server-side, via `GEMINI_API_KEY`)

## Quick start

### Prerequisites

- Node.js 18+
- Go 1.21+
- PostgreSQL 15+

### Backend (API)

```bash
cd services/api
cp .env.example .env   # set DATABASE_URL, JWT_SECRET, GEMINI_API_KEY
go run ./cmd/api
```

API runs at `http://localhost:8080` by default.

### Database

Run migrations once (use a fresh database or run the SQL manually; `type` is quoted for PostgreSQL compatibility).

**Option A — use the same URL as the API (e.g. Neon):**  
Set `DATABASE_URL` in your shell, then run psql. The `.env` in `services/api` is not loaded by the shell, so export it first:

```bash
cd infra/db
export DATABASE_URL="postgres://user:password@your-host.neon.tech/neondb?sslmode=require"   # or: source ../services/api/.env  if you have one
psql "$DATABASE_URL" -f migrations/001_schema.sql
```

**Option B — pass the URL directly:**

```bash
cd infra/db
psql "postgres://user:password@your-host.neon.tech/neondb?sslmode=require" -f migrations/001_schema.sql
```

Replace the URL with your real Neon (or other) connection string from `services/api/.env`. Use quotes so special characters in the password don’t break the command.

### Frontend (Web)

```bash
cd apps/web
cp .env.local.example .env.local   # set VITE_API_URL if needed
npm install
npm run dev
```

Web runs at `http://localhost:5173` by default.

## Project layout

```
/
  apps/web/           # Vite + React + shadcn/ui
  services/api/       # Go API (handler / service / repository / model / database)
  infra/db/           # SQL migrations
  plan.md             # Implementation plan
  prd.md              # Product requirements
```

## Environment

| Variable        | Where   | Description                    |
|----------------|---------|--------------------------------|
| `DATABASE_URL` | API     | PostgreSQL connection string   |
| `JWT_SECRET`   | API     | Secret for JWT signing         |
| `GEMINI_API_KEY` | API  | Gemini API key (content gen)    |
| `PORT`         | API     | HTTP port (default 8080)       |
| `CORS_ORIGINS` | API     | Allowed origins (default *)    |
| `VITE_API_URL` | Web     | API base URL (default http://localhost:8080) |
# Devflow_Os
