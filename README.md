# Disky

Disky is a price-tracking app for storage and memory products. The repository contains a Go API, a React frontend, PostgreSQL schema/migrations, and Docker Compose infrastructure for running the stack behind Caddy.

## Stack

- **Frontend:** React, TypeScript, Vite
- **Backend:** Go, Gin, pgx
- **Database:** PostgreSQL
- **Infrastructure:** Docker Compose, Caddy, Nginx

## Current capabilities

- Browse storage and RAM products with price-per-GB data
- Filter between storage and memory items in the frontend
- Serve product data from PostgreSQL, with mocked fallback data for empty development databases
- Stub OAuth entry points for Google, Apple, and Microsoft
- Stub alert endpoints for price alerts

## Repository layout

- `/frontend` – Vite/React application
- `/backend` – Go API server
- `/database/migrations` – PostgreSQL schema and seed data
- `/caddy` – reverse proxy configuration
- `/docker-compose.yml` – local stack definition
- `/.env.example` – environment variable template

## Getting started

1. Copy the environment template:

   ```bash
   cp .env.example .env
   ```

2. Update the values in `.env`, especially:
   - `JWT_SECRET`
   - PostgreSQL credentials
   - Amazon API credentials
   - OAuth client IDs and secrets
   - SMTP settings

3. Start the stack:

   ```bash
   docker compose up -d
   ```

The PostgreSQL container loads `/database/migrations/001_initial_schema.sql` on first startup.

## Local development

### Backend

```bash
cd backend
go mod download
go run ./cmd/server
```

The API listens on `SERVER_PORT` and exposes a health check at `/health`.

### Frontend

```bash
cd frontend
npm ci
npm run dev
```

For production builds, the frontend uses `VITE_API_URL` and `VITE_GOOGLE_CLIENT_ID` from the environment.

## API endpoints

- `GET /health`
- `GET /api/products`
- `GET /api/auth/google`
- `GET /api/auth/apple`
- `GET /api/auth/microsoft`
- `GET /api/alerts`
- `POST /api/alerts`

## Notes

- `docker-compose.yml` is configured to use published backend and frontend images from GHCR.
- The checked-in Caddy configuration is set up for `disky.tsew.com`.
- When the products table is empty, the API returns mock products for development/demo use.
