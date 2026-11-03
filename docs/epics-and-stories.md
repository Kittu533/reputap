# REPUTAP — Epics & Stories

**Based on:** Technical Design Document v1.0
**Development approach:** Phase-by-phase, backend-first per module

---

## EPIC 1: Foundation & Infrastructure

> Setup project skeleton, database, Docker, dan authentication sehingga project bisa jalan lokal.

### Stories

#### E1-S01: Project Setup & Docker Compose
**As** a developer
**I want** the project to run locally with Docker Compose
**So that** I can start developing immediately

**Tasks:**
- [ ] Bersihkan `go.mod` dari dependency lama (Gin, GORM)
- [ ] Setup Chi router + SQLX + golang-migrate dependencies
- [ ] Buat `cmd/server/main.go` entry point
- [ ] Buat `internal/config/config.go` (env loader)
- [ ] Buat `internal/database/mysql.go` (SQLX connection)
- [ ] Buat `docker-compose.yml` (MySQL 8 + Go app)
- [ ] Buat `.env.example`
- [ ] Buat `Makefile` (run, build, migrate-up, migrate-down, seed)

**Acceptance Criteria:**
- `docker compose up` menjalankan MySQL
- `make run` menjalankan Go server
- Server merespon `GET /health` dengan status 200

---

#### E1-S02: Database Migrations
**As** a developer
**I want** database schema managed via migrations
**So that** schema changes are tracked and reproducible

**Tasks:**
- [ ] Setup golang-migrate
- [ ] Migration 000001: `roles`, `permissions`, `role_permissions`
- [ ] Migration 000002: `users`
- [ ] Migration 000003: `customers`
- [ ] Migration 000004: `stores`
- [ ] Migration 000005: `products`
- [ ] Migration 000006: `orders`, `order_items`, `order_status_history`
- [ ] Migration 000007: `product_instances`, `qr_codes`
- [ ] Migration 000008: `scan_logs`
- [ ] Migration 000009: `design_assets`
- [ ] Migration 000010: `production_jobs`
- [ ] Migration 000011: `quality_controls`
- [ ] Migration 000012: `shipments`
- [ ] Migration 000013: `audit_logs`
- [ ] Buat seed data (roles, permissions, default super admin)

**Acceptance Criteria:**
- `make migrate-up` creates all tables
- `make migrate-down` drops all tables cleanly
- `make seed` inserts default roles + super admin account

---

#### E1-S03: API Foundation (Middleware & Response)
**As** a developer
**I want** standardized middleware and response helpers
**So that** all modules have consistent behavior

**Tasks:**
- [ ] Buat `internal/response/response.go` (Success, Error, Paginated)
- [ ] Buat `internal/middleware/logger.go` (request logging via slog)
- [ ] Buat `internal/middleware/cors.go`
- [ ] Buat `internal/middleware/recoverer.go` (panic recovery)
- [ ] Buat `pkg/pagination/pagination.go`
- [ ] Setup router di `main.go` dengan semua middleware

**Acceptance Criteria:**
- All responses use consistent JSON format
- Request logs show method, path, status, duration
- Panics are caught and return 500

---

#### E1-S04: Authentication System
**As** an admin
**I want** to login with email and password
**So that** I can access the dashboard securely

**Tasks:**
- [ ] Buat `internal/auth/model.go` (LoginRequest, TokenResponse, Claims)
- [ ] Buat `internal/auth/repository.go` (FindUserByEmail)
- [ ] Buat `internal/auth/service.go` (Login, GenerateTokens, RefreshToken)
- [ ] Buat `internal/auth/handler.go` (POST /login, POST /logout, GET /me, POST /refresh)
- [ ] Buat `internal/middleware/auth.go` (JWT validation middleware)
- [ ] Implement bcrypt password hashing
- [ ] Implement JWT access + refresh token

**Acceptance Criteria:**
- POST /api/v1/auth/login with valid creds returns tokens
- POST /api/v1/auth/login with invalid creds returns 401
- GET /api/v1/auth/me with valid token returns user info
- Protected endpoints reject requests without valid token
- Token refresh works correctly

---

## EPIC 2: Core Business Management

> CRUD untuk Customer, Store, Product, Order — supaya tim bisa mulai mencatat dan mengelola pesanan.

### Stories

#### E2-S01: Customer Management (Backend)
**As** an admin
**I want** to manage customer records
**So that** I can track who orders REPUTAP products

**Tasks:**
- [ ] Buat `internal/customers/model.go`
- [ ] Buat `internal/customers/repository.go` (Create, GetByID, List, Update, Delete, Search)
- [ ] Buat `internal/customers/service.go` (validation, customer_code generation)
- [ ] Buat `internal/customers/handler.go` (CRUD endpoints)
- [ ] Register routes di main router

**Acceptance Criteria:**
- CRUD operations work for all fields
- customer_code auto-generated (e.g., CUST-001)
- Search by name, phone, email, customer_code
- Pagination works
- Validation errors return proper messages

---

#### E2-S02: Store Management (Backend)
**As** an admin
**I want** to manage store/business records
**So that** each business has its Google Review URL configured

**Tasks:**
- [ ] Buat `internal/stores/model.go`
- [ ] Buat `internal/stores/repository.go`
- [ ] Buat `internal/stores/service.go` (URL validation, slug generation)
- [ ] Buat `internal/stores/handler.go`
- [ ] Implement Google Review URL validation (domain allowlist)
- [ ] Implement store status management

**Acceptance Criteria:**
- CRUD operations work
- Review URL validated against Google domain allowlist
- Store status can be changed (active/inactive/suspended)
- Slug auto-generated from store name
- Customer-store relationship works (1 customer → N stores)

---

#### E2-S03: Product Catalog (Backend)
**As** an admin
**I want** to manage REPUTAP product catalog
**So that** orders can reference specific products

**Tasks:**
- [ ] Buat `internal/products/model.go`
- [ ] Buat `internal/products/repository.go`
- [ ] Buat `internal/products/service.go`
- [ ] Buat `internal/products/handler.go`

**Acceptance Criteria:**
- CRUD operations work for products
- SKU must be unique
- Price fields (production_cost, selling_price) stored correctly
- Products can be filtered by status

---

#### E2-S04: Order Management (Backend)
**As** an admin
**I want** to create and manage orders
**So that** I can track the entire fulfillment process

**Tasks:**
- [ ] Buat `internal/orders/model.go` (Order, OrderItem, StatusHistory)
- [ ] Buat `internal/orders/repository.go`
- [ ] Buat `internal/orders/service.go` (order creation, status transitions, validation)
- [ ] Buat `internal/orders/handler.go`
- [ ] Implement order_number auto-generation (ORD-YYYYMM-XXXX)
- [ ] Implement status transition rules
- [ ] Implement order_status_history logging

**Acceptance Criteria:**
- Create order with customer, store, and order items
- Order number auto-generated
- Status transitions follow allowed paths
- Invalid transitions rejected with error
- Status history recorded with user and timestamp
- Order detail includes customer, store, items info
- Pagination and filtering work

**Status Transition Rules:**
```
pending → confirmed, cancelled, on_hold
confirmed → processing, cancelled, on_hold
processing → production, cancelled, on_hold
production → qc, on_hold
qc → ready_to_ship, production (rework), on_hold
ready_to_ship → shipped, on_hold
shipped → completed
on_hold → (previous status)
```

---

#### E2-S05: Customer Management (Frontend)
**As** an admin
**I want** a UI to manage customers
**So that** I don't have to use API directly

**Tasks:**
- [ ] Buat composable `useApi` (fetch wrapper with auth)
- [ ] Buat layout `admin.vue` (sidebar navigation)
- [ ] Buat page `pages/customers/index.vue` (list + search)
- [ ] Buat page `pages/customers/[id].vue` (detail)
- [ ] Buat component `CustomerForm.vue` (create/edit modal)
- [ ] Implement auth middleware (redirect to login if no token)

---

#### E2-S06: Store Management (Frontend)
- [ ] Buat page `pages/stores/index.vue`
- [ ] Buat page `pages/stores/[id].vue`
- [ ] Buat component `StoreForm.vue`
- [ ] Implement URL verification UI

---

#### E2-S07: Product Catalog (Frontend)
- [ ] Buat page `pages/products/index.vue`
- [ ] Buat page `pages/products/[id].vue`
- [ ] Buat component `ProductForm.vue`

---

#### E2-S08: Order Management (Frontend)
- [ ] Buat page `pages/orders/index.vue`
- [ ] Buat page `pages/orders/[id].vue` (detail with timeline)
- [ ] Buat component `OrderForm.vue` (create with customer/store/product selectors)
- [ ] Buat component `OrderStatusBadge.vue`
- [ ] Buat component `OrderTimeline.vue` (status history visualization)

---

## EPIC 3: QR & NFC Engine

> Sistem identifikasi produk, QR code generation, NFC URL, dan redirect engine — jantung dari REPUTAP.

### Stories

#### E3-S01: Identifier & Product Instance Generation (Backend)
**As** an admin
**I want** to generate product instances when order enters PROCESSING
**So that** each physical unit has a unique permanent identifier

**Tasks:**
- [ ] Buat `pkg/identifier/nanoid.go` (NanoID generator with RP- prefix)
- [ ] Buat `internal/productinstances/model.go`
- [ ] Buat `internal/productinstances/repository.go`
- [ ] Buat `internal/productinstances/service.go` (batch generate from order items)
- [ ] Buat `internal/productinstances/handler.go`

**Acceptance Criteria:**
- POST /api/v1/orders/{id}/instances/generate creates N instances based on quantity
- Each instance gets unique serial_number (RP-xxxxxxxx)
- Instances linked to order_item and store
- Collision detection (regenerate if duplicate)

---

#### E3-S02: QR Code Generation (Backend)
**As** an admin
**I want** QR codes generated for each product instance
**So that** the QR can be printed on the acrylic display

**Tasks:**
- [ ] Buat `internal/qrcodes/model.go`
- [ ] Buat `internal/qrcodes/repository.go`
- [ ] Buat `internal/qrcodes/service.go` (generate QR image, download)
- [ ] Buat `internal/qrcodes/handler.go`
- [ ] Integrate go-qrcode library
- [ ] QR code = URL `{REDIRECT_BASE_URL}/{serial_number}`

**Acceptance Criteria:**
- QR code auto-created when product instance is generated
- QR image downloadable as PNG and SVG
- QR encodes permanent URL: `https://reputap.id/r/RP-a7k9x2m4`
- QR status management (active/inactive/suspended/retired)

---

#### E3-S03: Redirect Engine (Backend)
**As** a store visitor
**I want** to be redirected to Google Review when I scan QR/tap NFC
**So that** I can leave a review easily

**Tasks:**
- [ ] Buat `internal/redirect/handler.go`
- [ ] Buat `internal/redirect/service.go` (resolve code, validate, redirect)
- [ ] Buat `internal/redirect/repository.go` (find QR + store by code)
- [ ] Buat `internal/redirect/templates/not_found.html`
- [ ] Buat `internal/redirect/templates/inactive.html`
- [ ] Implement async scan logging (goroutine)
- [ ] Implement device type detection from User-Agent
- [ ] Implement 302 redirect with no-cache headers

**Acceptance Criteria:**
- GET /r/RP-a7k9x2m4 → 302 to Google Review URL
- GET /r/RP-a7k9x2m4?src=qr → logged as QR scan
- GET /r/RP-a7k9x2m4?src=nfc → logged as NFC tap
- Invalid code → HTML error page
- Inactive QR → HTML inactive page
- Inactive store → HTML inactive page
- P95 response time < 100ms (local)

---

#### E3-S04: NFC Configuration Tracking (Backend)
**As** an operator
**I want** to update NFC status for each product instance
**So that** the team knows which units have NFC configured

**Tasks:**
- [ ] Add NFC status endpoints to productinstances handler
- [ ] PATCH /api/v1/product-instances/{id}/nfc-status
- [ ] Track nfc_configured_at and nfc_configured_by

**Acceptance Criteria:**
- NFC status can be updated through API
- NFC URL displayed: `{REDIRECT_BASE_URL}/{serial}?src=nfc`
- Status transitions: not_configured → configured → tested

---

#### E3-S05: QR Management (Frontend)
- [ ] Buat page `pages/qr-codes/index.vue` (list all QR codes)
- [ ] Buat component `QRPreview.vue` (show QR image + info)
- [ ] Buat download button (PNG/SVG)
- [ ] Show NFC URL for each instance

---

#### E3-S06: Instance Generation (Frontend)
- [ ] Add "Generate Instances" button on order detail page
- [ ] Show instances list on order detail
- [ ] Show NFC status per instance
- [ ] Show QR code preview per instance

---

## EPIC 4: Internal Operations

> Design management, production workflow, QC, shipping, dan audit log.

### Stories

#### E4-S01: Design Asset Management (Backend)
**Tasks:**
- [ ] Buat `internal/design/` module (model, repo, service, handler)
- [ ] File upload endpoint (PNG, JPG, SVG, PDF)
- [ ] File size validation (max 10MB)
- [ ] File type validation
- [ ] Version tracking
- [ ] Approval status management
- [ ] Store uploaded files to local `./uploads/designs/`

---

#### E4-S02: Production Management (Backend)
**Tasks:**
- [ ] Buat `internal/production/` module
- [ ] Production job CRUD
- [ ] Operator assignment
- [ ] Checklist data (JSON field)
- [ ] Status updates with timestamps
- [ ] Production notes

---

#### E4-S03: Quality Control (Backend)
**Tasks:**
- [ ] Buat `internal/qc/` module
- [ ] QC submission with checklist (JSON)
- [ ] Evidence upload (photos)
- [ ] Result: passed/failed/rework_required
- [ ] Update product_instance.qc_status
- [ ] Block: failed QC → cannot set order to ready_to_ship

---

#### E4-S04: Shipping Management (Backend)
**Tasks:**
- [ ] Buat `internal/shipping/` module
- [ ] Shipment CRUD
- [ ] Tracking number input
- [ ] Status management (waiting → packed → shipped → delivered)

---

#### E4-S05: Audit Log (Backend)
**Tasks:**
- [ ] Buat `internal/audit/` module
- [ ] Audit log recording service (reusable across modules)
- [ ] Log list endpoint with filtering (entity_type, date range)
- [ ] Immutable: no update/delete endpoints

---

#### E4-S06: Operations Frontend
- [ ] Design management page
- [ ] Production queue page with checklist
- [ ] QC submission page with photo upload
- [ ] Shipping page with tracking
- [ ] Activity log page

---

## EPIC 5: Dashboard & Analytics

> Dashboard overview, scan analytics, dan reporting.

### Stories

#### E5-S01: Dashboard Summary API (Backend)
**Tasks:**
- [ ] GET /api/v1/dashboard/summary
  - Total orders (by status)
  - Total customers
  - Total stores
  - Active QR codes
  - Total scans (QR vs NFC)
  - Redirect success rate
- [ ] GET /api/v1/dashboard/recent-activities (from audit_logs)

---

#### E5-S02: Scan Analytics API (Backend)
**Tasks:**
- [ ] GET /api/v1/analytics/overview (total, QR vs NFC, device breakdown)
- [ ] GET /api/v1/analytics/scans (time-series: daily/weekly/monthly)
- [ ] GET /api/v1/analytics/stores/{id} (per-store analytics)

---

#### E5-S03: Dashboard & Analytics (Frontend)
- [ ] Overview page with summary cards
- [ ] Order chart (by status, by time)
- [ ] Scan chart (daily, QR vs NFC)
- [ ] Recent activities feed
- [ ] Analytics page with store comparison

---

## EPIC 6: Login & Settings Frontend

### Stories

#### E6-S01: Login Page (Frontend)
- [ ] Buat page `pages/login.vue`
- [ ] Buat composable `useAuth` (login, logout, token management)
- [ ] Buat middleware `auth.ts` (route guard)
- [ ] Token storage in localStorage/cookie
- [ ] Auto-refresh token before expiry

---

#### E6-S02: Settings Page (Frontend)
- [ ] User profile / change password
- [ ] User management (Super Admin only)
- [ ] Role/permission management (Super Admin only)

---

## DEVELOPMENT PRIORITY ORDER

```
Phase 1: E1-S01 → E1-S02 → E1-S03 → E1-S04 → E6-S01
Phase 2: E2-S01 → E2-S02 → E2-S03 → E2-S04 → E2-S05 → E2-S06 → E2-S07 → E2-S08
Phase 3: E3-S01 → E3-S02 → E3-S03 → E3-S04 → E3-S05 → E3-S06
Phase 4: E4-S01 → E4-S02 → E4-S03 → E4-S04 → E4-S05 → E4-S06
Phase 5: E5-S01 → E5-S02 → E5-S03
```
