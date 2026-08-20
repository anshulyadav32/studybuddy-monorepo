# Architecture overview

This project follows a staged monorepo pattern for StudyBuddy:

- Next.js apps for student, admin, and teacher experiences
- Go API gateway and core services for business logic and routing
- Rust services for media processing and learning analytics
- PostgreSQL + Prisma as the system of record
- Redis as session/cache layer
- MinIO/S3-compatible storage for media assets

The service boundaries are intentionally modular so the platform can grow without large coupling between learning features, user management, and analytics.
