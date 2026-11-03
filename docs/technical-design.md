# REPUTAP — Technical Design Document

**Version:** 1.0
**Based on:** PRD v2.0
**Last Updated:** 2026-11-04

---

## 1. ARCHITECTURAL DECISIONS

Keputusan teknis hasil grilling session PRD.

### 1.1 Finalized Decisions

| # | Decision | Choice | Rationale |
|---|---|---|---|
| D-01 | review_url scope | 1 store = 1 review URL | MVP simplicity. Multi-destination = future. |
| D-02 | Circular FK resolution | Hapus `qr_code_id` di product_instances | QR tahu instance-nya via `product_instance_id`. Satu arah, no circular. |
| D-03 | Order-Store relation | 1 order = 1 store | Simple. Multi-store order = future. |
| D-04 | Identifier format | NanoID: `RP-` + 8 chars lowercase alphanumeric | Unpredictable, collision-safe, URL-friendly. Contoh: `RP-a7k9x2m4` |
| D-05 | Identifier = serial + QR | `serial_number` = `qr_codes.code` = sama | 1 identitas permanen per unit fisik. |
| D-06 | HTTP redirect | 302 + `Cache-Control: no-store, no-cache` | Destination bisa berubah, jangan di-cache browser. |
| D-07 | Scan logging | Async goroutine (fire-and-forget) | Non-blocking redirect. Acceptable data loss risk untuk MVP. |
| D-08 | Fallback error page | Go serve HTML statis | Tanpa Nuxt hop. Ringan. Self-contained. |
| D-09 | Order status | Simplified (8 status) | 14 terlalu banyak untuk MVP 5-10 order. |
| D-10 | Status transition | Semua manual | Admin klik untuk transisi. Predictable. |
| D-11 | Product instance creation | Saat admin trigger di status CONFIGURATION | Generate instances + QR + NFC URL sekaligus. |
| D-12 | Production job scope | Per product_instance | Granular. Tiap unit fisik bisa punya masalah berbeda. |
| D-13 | NFC API path | `/api/v1/product-instances/{id}/nfc` | NFC per instance, bukan per product catalog. |
| D-14 | Auth mechanism | JWT (access + refresh token) | Stateless, scalable. |
| D-15 | Backend framework | Chi Router + SQLX + golang-migrate | Go-idiomatic. Direct SQL control. Best untuk belajar. |
| D-16 | Permission model | Dynamic (roles + permissions + role_permissions tables) | Flexible, extensible. |
| D-17 | Order status history | Tabel `order_status_history` terpisah | Queryable, clean. Audit logs untuk hal lain. |
| D-18 | Google Place ID | Tidak perlu untuk MVP | Cukup URL saja. |

### 1.2 Simplified Order Status (8 status)

**Main flow:**

```
PENDING → CONFIRMED → PROCESSING → PRODUCTION → QC → READY_TO_SHIP → SHIPPED → COMPLETED
```

**Additional:**
- `CANCELLED`
- `ON_HOLD`

**Mapping dari PRD 14 status → Simplified 8:**

| PRD Status | Simplified | Keterangan |
|---|---|---|
| NEW_ORDER | PENDING | Order baru masuk |
| WAITING_CUSTOMER_DATA | PENDING | Masih dikategorikan pending |
| DATA_VERIFICATION | CONFIRMED | Data sudah diverifikasi |
| CONFIGURATION | PROCESSING | QR/NFC/Design sedang dikerjakan |
| DESIGN | PROCESSING | Termasuk dalam processing |
| DESIGN_APPROVAL | PROCESSING | Approval via WA, tracked di notes |
| PRODUCTION | PRODUCTION | Produksi fisik |
| QUALITY_CONTROL | QC | Quality control |
| READY_TO_SHIP | READY_TO_SHIP | Siap kirim |
| SHIPPED | SHIPPED | Dalam pengiriman |
| COMPLETED | COMPLETED | Selesai |

### 1.3 Google Domain Allowlist

```go
var AllowedGoogleDomains = []string{
    "search.google.com",
    "www.google.com",
    "google.com",
    "maps.google.com",
    "www.google.co.id",
    "google.co.id",
    "maps.app.goo.gl",
    "goo.gl",
    "g.page",
    "g.co",
}
```

---

## 2. TECHNOLOGY STACK

### 2.1 Backend

| Component | Technology | Version |
|---|---|---|
| Language | Go | 1.27+ |
| HTTP Router | chi | v5 |
| Database Driver | go-sql-driver/mysql | latest |
| SQL Toolkit | sqlx | latest |
| Migration | golang-migrate | v4 |
| Auth | golang-jwt/jwt | v5 |
| Password Hash | golang.org/x/crypto (bcrypt) | latest |
| Validation | go-playground/validator | v10 |
| QR Generation | skip2/go-qrcode | latest |
| Logging | log/slog (stdlib) | built-in |
| Env Config | joho/godotenv | latest |
| ID Generation | matoous/go-nanoid | v2 |
| CORS | go-chi/cors | latest |

### 2.2 Frontend

| Component | Technology | Version |
|---|---|---|
| Framework | Nuxt | 4 |
| UI | Vue | 3 |
| Language | TypeScript | 7+ |
| UI Components | Nuxt UI | 4 |
| CSS | Tailwind CSS | 4 |
| State | Pinia (jika perlu) | latest |

### 2.3 Infrastructure

| Component | Technology |
|---|---|
| Database | MySQL 8 |
| Containerization | Docker Compose |
| Package Manager (FE) | bun |

---

## 3. PROJECT STRUCTURE

```
reputap/
├── reputap-be/                    # Go Backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go            # Entry point
│   ├── internal/
│   │   ├── config/                # App configuration
│   │   │   └── config.go
│   │   ├── database/              # DB connection
│   │   │   └── mysql.go
│   │   ├── middleware/            # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   ├── response/             # Standardized API response
│   │   │   └── response.go
│   │   ├── auth/                  # Module: Authentication
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── users/                 # Module: User management
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── customers/             # Module: Customer management
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── stores/                # Module: Store management
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── orders/                # Module: Order management
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── products/              # Module: Product catalog
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── productinstances/      # Module: Physical unit tracking
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── qrcodes/               # Module: QR code management
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── redirect/              # Module: Redirect engine
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── templates/
│   │   │       ├── not_found.html
│   │   │       └── inactive.html
│   │   ├── analytics/             # Module: Scan analytics
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── production/            # Module: Production management
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── qc/                    # Module: Quality control
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   ├── shipping/              # Module: Shipping
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── model.go
│   │   └── audit/                 # Module: Audit logs
│   │       ├── handler.go
│   │       ├── service.go
│   │       ├── repository.go
│   │       └── model.go
│   ├── migrations/                # SQL migration files
│   │   ├── 000001_create_roles.up.sql
│   │   ├── 000001_create_roles.down.sql
│   │   └── ...
│   ├── pkg/                       # Shared utilities
│   │   ├── identifier/
│   │   │   └── nanoid.go
│   │   ├── validator/
│   │   │   └── url.go
│   │   └── pagination/
│   │       └── pagination.go
│   ├── go.mod
│   ├── go.sum
│   ├── .env.example
│   └── Makefile
│
├── reputap-fe/                    # Nuxt Frontend
│   ├── app/
│   │   ├── pages/
│   │   ├── components/
│   │   ├── layouts/
│   │   ├── composables/
│   │   ├── middleware/
│   │   └── types/
│   ├── public/
│   ├── nuxt.config.ts
│   └── package.json
│
├── docker-compose.yml
├── docs/
│   ├── technical-design.md        # This file
│   ├── epics-and-stories.md
│   └── todo.md
│
└── prd-REPUTAP.md
```

### 3.1 Module Pattern (Backend)

Setiap module mengikuti pattern yang sama:

```
module/
├── model.go       # Struct definitions (domain models, DTOs)
├── repository.go  # Database queries (SQLX)
├── service.go     # Business logic
└── handler.go     # HTTP handlers (request/response)
```

**Layer responsibilities:**

- **model.go** — Data structures. Tidak boleh import package lain dari internal.
- **repository.go** — SQL queries. Hanya bicara dengan database via SQLX.
- **service.go** — Business logic. Orchestrate repository calls. Validation rules.
- **handler.go** — Parse HTTP request, call service, format HTTP response.

**Dependency flow:**
```
Handler → Service → Repository → Database
```

---

## 4. DATABASE SCHEMA (Revised)

### 4.1 Entity Relationship

```
roles 1──N users
users 1──N orders (created_by)
users 1──N production_jobs (assigned_to)
users 1──N quality_controls (inspected_by)
users 1──N design_assets (uploaded_by)

customers 1──N stores
customers 1──N orders

stores 1──N orders
stores 1──N product_instances
stores 1──N qr_codes

orders 1──N order_items
order_items 1──N product_instances

products 1──N order_items

product_instances 1──1 qr_codes (via qr_codes.product_instance_id)
product_instances 1──N production_jobs
product_instances 1──N quality_controls

qr_codes 1──N scan_logs

orders 1──N order_status_history
orders 1──1 shipments
orders 1──N design_assets
```

### 4.2 Table Definitions

#### roles
```sql
CREATE TABLE roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### permissions
```sql
CREATE TABLE permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    module VARCHAR(50) NOT NULL
);
```

#### role_permissions
```sql
CREATE TABLE role_permissions (
    role_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);
```

#### users
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);
```

#### customers
```sql
CREATE TABLE customers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    customer_code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    business_name VARCHAR(200),
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_customer_code (customer_code),
    INDEX idx_phone (phone),
    INDEX idx_email (email)
);
```

#### stores
```sql
CREATE TABLE stores (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT UNSIGNED NOT NULL,
    store_code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    category ENUM('coffee_shop', 'restaurant', 'barbershop', 'salon', 'clinic', 'retail', 'hotel', 'other') DEFAULT 'other',
    address TEXT,
    google_business_name VARCHAR(255),
    review_url VARCHAR(500),
    review_url_verified BOOLEAN DEFAULT FALSE,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    INDEX idx_store_code (store_code),
    INDEX idx_slug (slug),
    INDEX idx_status (status)
);
```

#### products
```sql
CREATE TABLE products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sku VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    material VARCHAR(100),
    width_mm DECIMAL(8,2),
    height_mm DECIMAL(8,2),
    thickness_mm DECIMAL(8,2),
    nfc_chip_type VARCHAR(50) DEFAULT 'NTAG213',
    finishing VARCHAR(50),
    production_cost DECIMAL(12,2) DEFAULT 0,
    selling_price DECIMAL(12,2) DEFAULT 0,
    status ENUM('active', 'inactive', 'discontinued') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sku (sku)
);
```

#### orders
```sql
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_number VARCHAR(30) NOT NULL UNIQUE,
    customer_id BIGINT UNSIGNED NOT NULL,
    store_id BIGINT UNSIGNED NOT NULL,
    order_status ENUM('pending', 'confirmed', 'processing', 'production', 'qc', 'ready_to_ship', 'shipped', 'completed', 'cancelled', 'on_hold') DEFAULT 'pending',
    payment_status ENUM('unpaid', 'partial', 'paid', 'refunded') DEFAULT 'unpaid',
    subtotal DECIMAL(12,2) DEFAULT 0,
    discount DECIMAL(12,2) DEFAULT 0,
    total DECIMAL(12,2) DEFAULT 0,
    notes TEXT,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_order_number (order_number),
    INDEX idx_order_status (order_status),
    INDEX idx_customer_id (customer_id),
    INDEX idx_created_at (created_at)
);
```

#### order_items
```sql
CREATE TABLE order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    quantity INT UNSIGNED NOT NULL DEFAULT 1,
    unit_price DECIMAL(12,2) NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    INDEX idx_order_id (order_id)
);
```

#### product_instances
```sql
CREATE TABLE product_instances (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_item_id BIGINT UNSIGNED NOT NULL,
    store_id BIGINT UNSIGNED NOT NULL,
    serial_number VARCHAR(20) NOT NULL UNIQUE,
    nfc_status ENUM('not_configured', 'configured', 'tested', 'failed', 'replaced') DEFAULT 'not_configured',
    nfc_configured_at TIMESTAMP NULL,
    nfc_configured_by BIGINT UNSIGNED NULL,
    production_status ENUM('pending', 'in_progress', 'completed', 'failed') DEFAULT 'pending',
    qc_status ENUM('pending', 'passed', 'failed', 'rework') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_item_id) REFERENCES order_items(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (nfc_configured_by) REFERENCES users(id),
    UNIQUE INDEX idx_serial_number (serial_number)
);
```

#### qr_codes
```sql
CREATE TABLE qr_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_instance_id BIGINT UNSIGNED NOT NULL UNIQUE,
    store_id BIGINT UNSIGNED NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    status ENUM('active', 'inactive', 'suspended', 'retired') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    UNIQUE INDEX idx_code (code)
);
```

**Note:** `code` = `serial_number` dari product_instances. Keduanya menyimpan value yang sama (`RP-a7k9x2m4`). `qr_codes.code` digunakan untuk lookup redirect, `product_instances.serial_number` untuk identifikasi fisik.

#### scan_logs
```sql
CREATE TABLE scan_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    qr_code_id BIGINT UNSIGNED NOT NULL,
    source ENUM('qr', 'nfc', 'unknown') DEFAULT 'unknown',
    device_type ENUM('android', 'ios', 'desktop', 'other') DEFAULT 'other',
    user_agent TEXT,
    referrer VARCHAR(500),
    scanned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (qr_code_id) REFERENCES qr_codes(id),
    INDEX idx_qr_code_id (qr_code_id),
    INDEX idx_scanned_at (scanned_at),
    INDEX idx_source (source)
);
```

#### order_status_history
```sql
CREATE TABLE order_status_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    from_status VARCHAR(30),
    to_status VARCHAR(30) NOT NULL,
    changed_by BIGINT UNSIGNED NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (changed_by) REFERENCES users(id),
    INDEX idx_order_id (order_id)
);
```

#### design_assets
```sql
CREATE TABLE design_assets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    store_id BIGINT UNSIGNED NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(10) NOT NULL,
    file_size BIGINT UNSIGNED DEFAULT 0,
    version INT UNSIGNED DEFAULT 1,
    approval_status ENUM('draft', 'pending_approval', 'approved', 'revision_required') DEFAULT 'draft',
    uploaded_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (uploaded_by) REFERENCES users(id),
    INDEX idx_order_id (order_id)
);
```

#### production_jobs
```sql
CREATE TABLE production_jobs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_instance_id BIGINT UNSIGNED NOT NULL,
    assigned_to BIGINT UNSIGNED NULL,
    status ENUM('pending', 'in_progress', 'completed', 'failed') DEFAULT 'pending',
    checklist_data JSON,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (assigned_to) REFERENCES users(id),
    INDEX idx_status (status),
    INDEX idx_assigned_to (assigned_to)
);
```

#### quality_controls
```sql
CREATE TABLE quality_controls (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_instance_id BIGINT UNSIGNED NOT NULL,
    inspected_by BIGINT UNSIGNED NOT NULL,
    result ENUM('passed', 'failed', 'rework_required') NOT NULL,
    checklist_data JSON NOT NULL,
    evidence_paths JSON,
    notes TEXT,
    inspected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (inspected_by) REFERENCES users(id),
    INDEX idx_product_instance_id (product_instance_id)
);
```

#### shipments
```sql
CREATE TABLE shipments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    recipient_name VARCHAR(150) NOT NULL,
    recipient_address TEXT NOT NULL,
    recipient_phone VARCHAR(20),
    courier VARCHAR(100),
    tracking_number VARCHAR(100),
    shipping_status ENUM('waiting', 'packed', 'shipped', 'delivered', 'failed') DEFAULT 'waiting',
    shipped_at TIMESTAMP NULL,
    delivered_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    INDEX idx_order_id (order_id),
    INDEX idx_shipping_status (shipping_status)
);
```

#### audit_logs
```sql
CREATE TABLE audit_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT UNSIGNED,
    old_values JSON,
    new_values JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_created_at (created_at)
);
```

---

## 5. API DESIGN

### 5.1 Standard Response Format

**Success:**
```json
{
    "success": true,
    "data": { ... },
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

**Error:**
```json
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Email is required",
        "details": [
            { "field": "email", "message": "Email is required" }
        ]
    }
}
```

### 5.2 API Endpoints (Revised)

**Auth**
```
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/me
POST   /api/v1/auth/refresh
```

**Customers**
```
GET    /api/v1/customers
POST   /api/v1/customers
GET    /api/v1/customers/{id}
PUT    /api/v1/customers/{id}
DELETE /api/v1/customers/{id}
```

**Stores**
```
GET    /api/v1/stores
POST   /api/v1/stores
GET    /api/v1/stores/{id}
PUT    /api/v1/stores/{id}
PATCH  /api/v1/stores/{id}/status
PUT    /api/v1/stores/{id}/review-url
```

**Products**
```
GET    /api/v1/products
POST   /api/v1/products
GET    /api/v1/products/{id}
PUT    /api/v1/products/{id}
```

**Orders**
```
GET    /api/v1/orders
POST   /api/v1/orders
GET    /api/v1/orders/{id}
PUT    /api/v1/orders/{id}
PATCH  /api/v1/orders/{id}/status
GET    /api/v1/orders/{id}/history
```

**Product Instances & QR**
```
GET    /api/v1/orders/{id}/instances            # List instances for an order
POST   /api/v1/orders/{id}/instances/generate    # Generate instances + QR
GET    /api/v1/product-instances/{id}
PATCH  /api/v1/product-instances/{id}/nfc-status
```

**QR Codes**
```
GET    /api/v1/qr-codes
GET    /api/v1/qr-codes/{id}
GET    /api/v1/qr-codes/{id}/download?format=png|svg
PATCH  /api/v1/qr-codes/{id}/status
```

**Redirect (Public)**
```
GET    /r/{code}
```

**Production**
```
GET    /api/v1/production/jobs
POST   /api/v1/production/jobs
GET    /api/v1/production/jobs/{id}
PATCH  /api/v1/production/jobs/{id}/status
```

**Quality Control**
```
POST   /api/v1/quality-controls
GET    /api/v1/quality-controls/instance/{instanceId}
```

**Shipping**
```
GET    /api/v1/shipments
POST   /api/v1/shipments
GET    /api/v1/shipments/{id}
PATCH  /api/v1/shipments/{id}/status
```

**Analytics**
```
GET    /api/v1/analytics/overview
GET    /api/v1/analytics/scans
GET    /api/v1/analytics/stores/{id}
```

**Audit Logs**
```
GET    /api/v1/audit-logs
```

**Dashboard**
```
GET    /api/v1/dashboard/summary
GET    /api/v1/dashboard/recent-activities
```

---

## 6. REDIRECT ENGINE FLOW

```
┌──────────────────────────────────────────────────────────┐
│  Visitor scans QR / taps NFC                             │
│  URL: https://reputap.id/r/RP-a7k9x2m4?src=qr           │
└────────────────────┬─────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────────┐
│  Go Backend: GET /r/{code}                               │
│  1. Extract code: RP-a7k9x2m4                            │
│  2. Extract source: qr (from ?src= param)                │
└────────────────────┬─────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────────┐
│  Query: qr_codes WHERE code = 'RP-a7k9x2m4'             │
│  JOIN stores ON qr_codes.store_id = stores.id            │
└────────────────────┬─────────────────────────────────────┘
                     │
              ┌──────┴──────┐
              │             │
         Found?          Not Found
              │             │
              ▼             ▼
     QR active?      Return 404 HTML
     Store active?    "QR tidak ditemukan"
     URL valid?
              │
       ┌──────┴──────┐
       │             │
      Yes           No
       │             │
       ▼             ▼
  Async: log      Return HTML
  scan event      "QR tidak aktif"
       │
       ▼
  HTTP 302
  Location: {store.review_url}
  Cache-Control: no-store, no-cache
```

---

## 7. AUTH FLOW (JWT)

```
┌─────────────────────────────────┐
│  POST /api/v1/auth/login        │
│  Body: { email, password }      │
└───────────────┬─────────────────┘
                │
                ▼
┌─────────────────────────────────┐
│  Verify credentials             │
│  Generate:                      │
│  - Access Token (15 min)        │
│  - Refresh Token (7 days)       │
└───────────────┬─────────────────┘
                │
                ▼
┌─────────────────────────────────┐
│  Response:                      │
│  { access_token, refresh_token, │
│    expires_in, user }           │
└─────────────────────────────────┘

Subsequent requests:
Authorization: Bearer {access_token}

Token refresh:
POST /api/v1/auth/refresh
Body: { refresh_token }
→ New access_token + refresh_token
```

---

## 8. IDENTIFIER GENERATION

### NanoID Configuration

```go
const (
    IdentifierPrefix  = "RP-"
    IdentifierLength  = 8
    IdentifierAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
)

// Output examples:
// RP-a7k9x2m4
// RP-3bfn8yqz
// RP-p2wk6t9j
```

- 36^8 = ~2.8 trillion combinations
- Collision probability negligible untuk MVP scale
- URL-friendly, case-insensitive
- Human-readable (no ambiguous chars like 0/O, 1/l)

### Code-level: serial_number = qr code

Saat generate product instances:
1. Generate NanoID → `RP-a7k9x2m4`
2. Set `product_instances.serial_number` = `RP-a7k9x2m4`
3. Create `qr_codes.code` = `RP-a7k9x2m4`
4. Generate QR image dari URL `https://reputap.id/r/RP-a7k9x2m4`

---

## 9. ENVIRONMENT VARIABLES

```env
# Server
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=3306
DB_NAME=reputap
DB_USER=reputap
DB_PASSWORD=secret

# JWT
JWT_SECRET=your-256-bit-secret-here
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# App
APP_BASE_URL=http://localhost:8080
REDIRECT_BASE_URL=http://localhost:8080/r

# File Upload
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=10485760
```
