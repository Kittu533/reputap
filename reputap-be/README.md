# Reputap Backend (`reputap-be`)

RESTful API backend for Reputap NFC/QR smart review management platform.

## Tech Stack

- **Language**: Go 1.24+
- **HTTP Router**: Chi (`go-chi/chi/v5`)
- **Database**: PostgreSQL 16
- **Database Driver & Query**: `github.com/lib/pq` + `github.com/jmoiron/sqlx`
- **Authentication**: JWT (`golang-jwt/jwt/v5`) & bcrypt (`golang.org/x/crypto/bcrypt`)
- **API Documentation**: Swagger (`swaggo/http-swagger`)

## Environment Configuration

Create or update `.env` in `reputap-be/`:

```env
# Server
APP_PORT=8081
APP_ENV=development

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5433
DB_NAME=reputap
DB_USER=reputap
DB_PASSWORD=secret
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-256-bit-secret-change-this-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# App & Redirect Engine
APP_BASE_URL=http://localhost:8081
REDIRECT_BASE_URL=http://localhost:8081/r

# Uploads
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=10485760
```

## Setup & Running

### 1. Database & Migrations
From root directory:
```bash
make db
make migrate
```

Verify tables and seeded users:
```bash
go run cmd/verify/main.go
```

### 2. Run API Server
```bash
go run cmd/server/main.go
```
API runs on `http://localhost:8081`.

### 3. Build Production Binary
```bash
go build -o bin/server cmd/server/main.go
```

## API Modules

| Module | Base Path | Description |
|---|---|---|
| Health Check | `/health` | Server status |
| Redirect Engine | `/r/{code}` | Public tap/scan resolver & logger |
| Auth | `/api/v1/auth` | Login, refresh token, user profile |
| Customers | `/api/v1/customers` | Customer CRUD & management |
| Stores | `/api/v1/stores` | Store profiles & review URLs |
| Products | `/api/v1/products` | NFC stand/card product master |
| Orders | `/api/v1/orders` | Customer order management & items |
| Instances | `/api/v1/product-instances` | Serialized item tracking & NFC configs |
| QR Codes | `/api/v1/qr-codes` | Code generation & status |
| Design Assets | `/api/v1/design-assets` | Branding uploads & approvals |
| Production | `/api/v1/production` | Production job workflow & checklists |
| Quality Control| `/api/v1/qc` | QC inspection logs & pass/fail status |
| Shipping | `/api/v1/shipping` | Courier tracking & dispatch |
| Analytics | `/api/v1/analytics` | Tap/scan metrics & device stats |
| Dashboard | `/api/v1/dashboard` | High-level metrics & recent activities |
| Audit Logs | `/api/v1/audit-logs` | User action tracking |
| Swagger Docs | `/swagger/index.html` | Interactive API documentation |

## Default Admin Credentials
- **Email**: `admin@reputap.id`
- **Password**: `admin123`
