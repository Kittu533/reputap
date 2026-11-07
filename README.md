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
make seed     # Seeds demo data (customers, stores, products, orders, scan analytics)
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

## Landing page SEO

Set `NUXT_PUBLIC_SITE_URL` to the public frontend origin (currently `https://reputap-alpha.vercel.app`), without `/login`. The landing page renders its canonical URL, social preview, and product catalog structured data on the server. `/sitemap.xml` lists the public homepage; admin and activation routes use `noindex, nofollow`.

After deployment, verify the site in Google Search Console and submit `/sitemap.xml`. Update this URL when moving to a custom domain. Run SEO helper checks with `cd reputap-fe && bun test tests/seo.test.ts`. After `npm run build`, run `node tests/seo-smoke.mjs` to check the rendered metadata and SEO endpoints on temporary port 3197.
