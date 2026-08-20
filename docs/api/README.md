# StudyBuddy API description

This API is served through the public gateway at `/api/v2` and proxied to the core Go service. All protected endpoints require a bearer token in the `Authorization` header.

## Base URLs

- Gateway: `http://localhost:8080`
- Core service: `http://localhost:8081`

## Health

### GET /health
Returns gateway health status.

Response:
```json
{
  "status": "ok",
  "service": "gateway"
}
```

## Authentication

### POST /api/v2/auth/register
Creates a new student or teacher account.

Request body:
```json
{
  "name": "Aarav Sharma",
  "email": "aarav@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "token": "jwt-token",
  "user": {
    "id": "user_123",
    "name": "Aarav Sharma",
    "email": "aarav@example.com",
    "profile": {}
  }
}
```

### POST /api/v2/auth/login
Authenticates an existing user.

Request body:
```json
{
  "email": "aarav@example.com",
  "password": "password123"
}
```

Response: same shape as register.

### GET /api/v2/auth/validate
Validates a bearer token.

Headers:
```http
Authorization: Bearer <token>
```

Response:
```json
{
  "valid": true,
  "user": {
    "id": "user_123",
    "name": "Aarav Sharma",
    "email": "aarav@example.com"
  }
}
```

### GET /api/v2/auth/me
Returns the authenticated user identity.

## User profile

### GET /api/v2/users/me
Returns the current authenticated user with profile data.

### GET /api/v2/users/profile
Returns the authenticated user's profile object.

### PUT /api/v2/users/profile
Updates the authenticated user's profile.

Request body:
```json
{
  "phone": "+91-9876543210",
  "classLevel": "Class 10",
  "board": "CBSE",
  "preferences": {
    "studyGoal": "board-exam",
    "difficulty": "medium"
  }
}
```

Response:
```json
{
  "phone": "+91-9876543210",
  "classLevel": "Class 10",
  "board": "CBSE",
  "preferences": {
    "studyGoal": "board-exam",
    "difficulty": "medium"
  }
}
```

## Course and learning content

### GET /api/v2/courses
Returns all courses.

### POST /api/v2/courses
Creates a new course.

Request body:
```json
{
  "title": "Chemistry Fundamentals",
  "subjectId": "chem-101",
  "teacherId": "teacher-1"
}
```

Response:
```json
{
  "id": "course_123",
  "title": "Chemistry Fundamentals",
  "slug": "chemistry-fundamentals",
  "subjectId": "chem-101",
  "teacherId": "teacher-1",
  "createdAt": "2026-08-20T06:41:52Z"
}
```

### GET /api/v2/courses/:id
Returns a single course.

### GET /api/v2/courses/chapters?courseId=<courseId>
Returns all chapters for a course.

### POST /api/v2/courses/chapters?courseId=<courseId>
Creates a chapter for a course.

Request body:
```json
{
  "title": "Cells and Tissues",
  "order": 1
}
```

Response:
```json
{
  "id": "chapter_123",
  "courseId": "course_123",
  "title": "Cells and Tissues",
  "order": 1
}
```

## Planned future routes

The following endpoints are planned for later phases and are part of the roadmap:

- `/api/v2/learning/*`
- `/api/v2/tests/*`
- `/api/v2/payments/*`
- `/api/v2/notifications/*`
- `/api/v2/ai/*`

This document reflects the implemented contract from the current milestone and will expand as the platform grows.
