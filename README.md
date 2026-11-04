# Reputap

Smart Review & NFC/QR Tap Management System.

## Architecture

- **`reputap-be`**: REST API built with Go (Chi router, SQLX, PostgreSQL 16).
- **`reputap-fe`**: Admin Dashboard built with Nuxt 4 (Vue 3, Nuxt UI, Tailwind CSS, Bun).
- **Database**: PostgreSQL 16 containerized via Docker on host port `5433`.

## Quick Start

### 1. Prerequisites
- Docker & Docker Compose
- Go 1.24+
- Bun (`curl -fsSL https://bun.sh/install | bash`)

### 2. Run Database & Migrations
```bash
make db       # Starts reputap-postgres on port 5433
make migrate  # Runs all migrations into PostgreSQL
```

### 3. Run Full-Stack Development
```bash
make dev      # Runs both backend (:8081) and frontend (:3000)
```

Or run services individually:
```bash
make dev-be   # Start Go API server (port 8081)
make dev-fe   # Start Nuxt frontend (port 3000)
```

## Service URLs & Ports

| Service | Address | Default Port | Notes |
|---|---|---|---|
| Frontend UI | `http://localhost:3000` | `3000` | Nuxt 4 dashboard |
| Backend API | `http://localhost:8081` | `8081` | REST API |
| Swagger Docs | `http://localhost:8081/swagger/index.html` | `8081` | API Documentation |
| PostgreSQL | `localhost:5433` | `5433` (container `5432`) | Port 5433 mapped to avoid host collision |

## Default Credentials
- **Email**: `admin@reputap.id`
- **Password**: `admin123`
- **Role**: `super_admin`
