# StudyBuddy

StudyBuddy is a BYJU'S-style edtech platform built as a staged monorepo with Next.js frontends, Go backend services, Rust performance services, Prisma/PostgreSQL, Redis, and object storage.

## Monorepo structure

- `frontend/` — student web app
- `admin/` — admin dashboard
- `teacher/` — teacher portal
- `mobile/` — Flutter app
- `backend/gateway/` — Go API gateway
- `backend/core/` — Go business services
- `backend/learning/` — Rust learning engine
- `backend/video/` — Rust video processing service
- `db/prisma/` — Prisma schema and migration setup
- `storage/` — media/service storage conventions
- `docs/` — architecture and API references

## Local development

1. Copy `.env.example` to `.env` and adjust values.
2. Start infrastructure services:
   ```bash
   docker compose up -d postgres redis minio ngrok
   ```
3. If you want a public tunnel, add your Ngrok auth token to `.env` and set `NGROK_DOMAIN` to a reserved domain. Without a token, the `ngrok` service stays idle and only exposes the local dashboard on `http://localhost:4040`.
4. Install workspace dependencies:
   ```bash
   pnpm install
   ```
5. Run the app suite:
   ```bash
   pnpm dev
   ```

## Phase roadmap

1. Foundation and monorepo structure
2. Infrastructure and database design
3. Authentication and user management
4. Learning content and student experience
5. Teacher/admin flows
6. AI tutoring and recommendations
7. Payments and security
8. Production deployment

## Primary stack

- Frontend: Next.js + TypeScript + Tailwind
- Backend: Go (gateway + core services)
- Performance services: Rust (learning/video processing)
- Data: PostgreSQL + Prisma
- Cache: Redis
- Storage: MinIO/S3-compatible object storage
- Deployment: Docker + CI/CD

## Ngrok public tunnel (local dev)

Public URLs (created during local compose run):

- https://convent-fantasize-amicably.ngrok-free.dev
- http://convent-fantasize-amicably.ngrok-free.dev

Forwards to: http://host.docker.internal:3030 (host port 3030)

Ngrok web UI (inspect traffic): http://localhost:4040

Security: the ngrok authtoken is stored in your local .env file (NGROK_AUTHTOKEN). Do NOT commit .env to version control. If you need to rotate or remove the token, edit or delete the value in .env.
