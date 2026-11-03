# REPUTAP — Development Todo

**Status Legend:** ⬜ Not Started | 🟡 In Progress | ✅ Done

---

## PHASE 1 — Foundation

### 1.1 Project Setup
- ✅ Bersihkan go.mod dari Gin/GORM dependencies
- ✅ Install dependencies: chi, sqlx, golang-migrate, jwt, bcrypt, godotenv, nanoid, go-qrcode, cors, validator
- ✅ Buat `cmd/server/main.go`
- ✅ Buat `internal/config/config.go`
- ✅ Buat `internal/database/mysql.go`
- ✅ Buat `docker-compose.yml` (MySQL 8)
- ✅ Buat `.env.example` + `.env`
- ✅ Buat `Makefile`
- ⬜ Test: `docker compose up` + `make run` + `GET /health` → 200

### 1.2 Database Migrations
- ⬜ Setup golang-migrate binary/tool
- ✅ Migration 000001: roles, permissions, role_permissions
- ✅ Migration 000002: users
- ✅ Migration 000003: customers
- ✅ Migration 000004: stores
- ✅ Migration 000005: products
- ✅ Migration 000006: orders, order_items, order_status_history
- ✅ Migration 000007: product_instances, qr_codes
- ✅ Migration 000008: scan_logs
- ✅ Migration 000009: design_assets
- ✅ Migration 000010: production_jobs
- ✅ Migration 000011: quality_controls
- ✅ Migration 000012: shipments
- ✅ Migration 000013: audit_logs
- ✅ Seed: default roles (super_admin, admin, operator)
- ✅ Seed: default permissions
- ✅ Seed: default super admin user
- ⬜ Test: `make migrate-up` dan `make migrate-down`

### 1.3 API Foundation
- ✅ `internal/response/response.go` — Success(), Error(), Paginated()
- ✅ `internal/middleware/logger.go` — slog request logging
- ✅ `internal/middleware/cors.go` — CORS config
- ✅ `internal/middleware/recoverer.go` — panic recovery
- ✅ `pkg/pagination/pagination.go` — page, per_page, offset helper

### 1.4 Authentication
- ✅ `internal/auth/model.go` — LoginRequest, TokenResponse, Claims
- ✅ `internal/auth/repository.go` — FindUserByEmail
- ✅ `internal/auth/service.go` — Login, GenerateTokens, RefreshToken, ValidateToken
- ✅ `internal/auth/handler.go` — POST /login, POST /logout, GET /me, POST /refresh
- ✅ `internal/auth/middleware.go` — JWT validation middleware + RequirePermission
- ⬜ Test: login → get token → access protected endpoint → refresh → logout

### 1.5 Login Page (Frontend)
- ⬜ `app/composables/useApi.ts` — fetch wrapper with auth headers
- ⬜ `app/composables/useAuth.ts` — login, logout, token management
- ⬜ `app/middleware/auth.ts` — route guard
- ⬜ `app/pages/login.vue`
- ⬜ `app/layouts/admin.vue` — sidebar navigation skeleton

---

## PHASE 2 — Core Business Management

### 2.1 Customer Module (Backend)
- ⬜ `internal/customers/model.go`
- ⬜ `internal/customers/repository.go` — Create, GetByID, List, Update, Delete, Search
- ⬜ `internal/customers/service.go` — validation, code generation (CUST-001)
- ⬜ `internal/customers/handler.go` — GET, POST, GET/:id, PUT/:id, DELETE/:id
- ⬜ Test: CRUD + search + pagination

### 2.2 Store Module (Backend)
- ⬜ `internal/stores/model.go`
- ⬜ `internal/stores/repository.go`
- ⬜ `internal/stores/service.go` — URL validation, slug generation
- ⬜ `internal/stores/handler.go`
- ⬜ `pkg/validator/url.go` — Google domain allowlist validator
- ⬜ Test: CRUD + URL validation + status change

### 2.3 Product Module (Backend)
- ⬜ `internal/products/model.go`
- ⬜ `internal/products/repository.go`
- ⬜ `internal/products/service.go`
- ⬜ `internal/products/handler.go`
- ⬜ Test: CRUD + SKU uniqueness

### 2.4 Order Module (Backend)
- ⬜ `internal/orders/model.go` — Order, OrderItem, StatusHistory, CreateOrderRequest
- ⬜ `internal/orders/repository.go` — full CRUD + status history
- ⬜ `internal/orders/service.go` — status transition rules, validation, order_number gen
- ⬜ `internal/orders/handler.go` — CRUD + PATCH status + GET history
- ⬜ Test: create order → transition through statuses → verify history

### 2.5 Frontend Pages
- ⬜ `app/pages/customers/index.vue` — list, search, create modal
- ⬜ `app/pages/customers/[id].vue` — detail, edit
- ⬜ `app/components/customers/CustomerForm.vue`
- ⬜ `app/pages/stores/index.vue`
- ⬜ `app/pages/stores/[id].vue`
- ⬜ `app/components/stores/StoreForm.vue`
- ⬜ `app/pages/products/index.vue`
- ⬜ `app/pages/products/[id].vue`
- ⬜ `app/components/products/ProductForm.vue`
- ⬜ `app/pages/orders/index.vue`
- ⬜ `app/pages/orders/[id].vue`
- ⬜ `app/components/orders/OrderForm.vue`
- ⬜ `app/components/orders/OrderStatusBadge.vue`
- ⬜ `app/components/orders/OrderTimeline.vue`

---

## PHASE 3 — QR & NFC Engine

### 3.1 Identifier Generation
- ⬜ `pkg/identifier/nanoid.go` — GenerateID() → "RP-xxxxxxxx"
- ⬜ Test: uniqueness, format, collision handling

### 3.2 Product Instance Module (Backend)
- ⬜ `internal/productinstances/model.go`
- ⬜ `internal/productinstances/repository.go`
- ⬜ `internal/productinstances/service.go` — batch generate from order
- ⬜ `internal/productinstances/handler.go`
- ⬜ Test: generate instances → verify serial numbers + QR codes created

### 3.3 QR Code Module (Backend)
- ⬜ `internal/qrcodes/model.go`
- ⬜ `internal/qrcodes/repository.go`
- ⬜ `internal/qrcodes/service.go` — generate QR image (PNG/SVG)
- ⬜ `internal/qrcodes/handler.go` — list, detail, download, status
- ⬜ Test: generate → download PNG → download SVG → verify URL encoded

### 3.4 Redirect Engine (Backend)
- ⬜ `internal/redirect/repository.go` — FindActiveQRWithStore(code)
- ⬜ `internal/redirect/service.go` — resolve, validate, detect device
- ⬜ `internal/redirect/handler.go` — GET /r/{code}
- ⬜ `internal/redirect/templates/not_found.html`
- ⬜ `internal/redirect/templates/inactive.html`
- ⬜ Async scan logging (goroutine)
- ⬜ Test: valid code → 302 | invalid → 404 HTML | inactive → HTML | source tracking

### 3.5 NFC Status Tracking
- ⬜ PATCH endpoint for NFC status
- ⬜ NFC URL generation display
- ⬜ Test: status transitions

### 3.6 Frontend
- ⬜ QR Management page (list, preview, download)
- ⬜ Instance generation on order detail
- ⬜ NFC status display per instance

---

## PHASE 4 — Internal Operations

### 4.1 Design Asset Management
- ⬜ Backend: upload, list, version, approval status
- ⬜ File validation (type + size)
- ⬜ Frontend: design management page

### 4.2 Production Management
- ⬜ Backend: job CRUD, assignment, checklist, notes
- ⬜ Frontend: production queue page

### 4.3 Quality Control
- ⬜ Backend: QC submission, checklist, evidence, result
- ⬜ Business rule: failed QC blocks ready_to_ship
- ⬜ Frontend: QC page

### 4.4 Shipping Management
- ⬜ Backend: shipment CRUD, tracking, status
- ⬜ Frontend: shipping page

### 4.5 Audit Log
- ⬜ Backend: reusable audit service, list endpoint
- ⬜ Frontend: activity log page

---

## PHASE 5 — Dashboard & Analytics

### 5.1 Dashboard API
- ⬜ GET /api/v1/dashboard/summary
- ⬜ GET /api/v1/dashboard/recent-activities

### 5.2 Analytics API
- ⬜ GET /api/v1/analytics/overview
- ⬜ GET /api/v1/analytics/scans (time-series)
- ⬜ GET /api/v1/analytics/stores/{id}

### 5.3 Frontend
- ⬜ Overview page (summary cards + charts)
- ⬜ Analytics page (scan charts, store comparison)

---

## CURRENT FOCUS

> **Start with Phase 1.1 — Project Setup**
