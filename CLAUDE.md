# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Backend (Go)

```bash
cd backend
air          # dev server with hot reload
go run .     # run without hot reload
go build .   # build binary
go vet ./... # lint
```

### Frontend (Next.js)

```bash
cd frontend
pnpm install
pnpm dev     # dev server at http://localhost:3000
pnpm build   # production build
pnpm lint    # ESLint
```

## Environment Setup

**Backend** (`backend/.env`):
```
JWT_SECRET=your-secret
PORT=8080
GIN_MODE=debug
CLIENT_URL=http://localhost:3000
```

**Frontend** (`frontend/.env`):
```
JWT_SECRET=your-secret        # must match backend
NEXT_PUBLIC_BASE_URL=http://localhost:3000
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

`JWT_SECRET` must be identical on both sides. Generate with: `openssl rand -base64 32`

## Architecture

Monorepo with two independent services:

- `backend/` — Go REST API (Gin + GORM + SQLite)
- `frontend/` — Next.js 16 app (Auth.js v5 + Tailwind CSS v4)

### Auth Flow

1. Frontend `src/proxy.ts` (Next.js middleware) guards all routes — unauthenticated users redirect to `/sign-in`; authenticated users are restricted to `/dashboard` and `/profile`
2. Auth.js `CredentialsProvider` (`src/lib/auth/index.ts`) calls Go `POST /api/auth/login` and stores `accessToken`, `refreshToken`, and `accessTokenExpires` in the JWT session
3. On each request the JWT callback checks expiry and calls `POST /api/auth/refresh` automatically
4. If refresh fails, `session.error` is set to `"RefreshAccessTokenError"` — `AuthGuard` component detects this and calls `signOut`
5. Backend `middleware/auth.go` validates Bearer tokens on protected routes and sets `userID` / `email` in Gin context via `config.JWT_CLAIMS_KEY_*`

### API Client

`src/lib/api/index.ts` exports `BaseApi` — a thin wrapper around `ky`. All service files (`src/services/`) use `BaseApi`. Pass a `session` object to inject `Authorization: Bearer <token>` automatically. On 401, the `afterResponse` hook redirects to `/sign-in`.

### Backend Package Layout

| Package | Purpose |
|---|---|
| `config/` | DB connection/migration + env-sourced constants |
| `handlers/` | HTTP handler structs (one per domain) |
| `middleware/` | JWT auth middleware |
| `models/` | GORM model structs (auto-migrated on startup) |
| `enums/` | Typed enums (Role: Client/Admin/SuperAdmin) |
| `utils/` | Env helpers |

### Frontend Path Aliases

`@/` maps to `src/`. Component organisation:
- `src/components/core/` — app-wide structural components (e.g. `AuthGuard`)
- `src/components/common/` — shared UI components
- `src/components/pages/` — page-scoped components
- `src/lib/auth/` — Auth.js config and callbacks
- `src/lib/api/` — `BaseApi` HTTP client
- `src/services/` — typed API call definitions
- `src/types/core/next-auth.d.ts` — module augmentation extending `Session`, `User`, and `JWT` with `accessToken`, `refreshToken`, `role`

### Database

SQLite (`backend/database.db`, git-ignored). Schema auto-migrated via `config.MigrateDatabase` on startup. No migration files — GORM `AutoMigrate` handles schema evolution.
