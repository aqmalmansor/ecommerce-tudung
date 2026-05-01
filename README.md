# Next.js + Go JWT Authentication Example

A simple example demonstrating JWT authentication between a Go backend and Next.js frontend using Auth.js.

## Tech Stack

- **Frontend**: Next.js with Auth.js (NextAuth.js)
- **Backend**: Go with JWT

## Getting Started

### Backend

```bash
cd backend
air
```

### Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) to view the app.

## How It Works

1. User authenticates via the Next.js frontend
2. Auth.js handles the session and sends credentials to the Go backend
3. Go backend validates credentials and returns a JWT
4. JWT is stored in the session and used for authenticated API requests
