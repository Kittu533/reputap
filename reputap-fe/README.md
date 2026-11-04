# Reputap Frontend (`reputap-fe`)

Admin Dashboard interface for Reputap NFC/QR Smart Review Platform.

## Tech Stack

- **Framework**: Nuxt 4 (Vue 3, Composition API)
- **UI Components**: Nuxt UI (`@nuxt/ui`) + Tailwind CSS
- **Icons**: Lucide Icons (`@iconify-json/lucide`)
- **Package Manager & Runtime**: Bun

## Environment Configuration

Configure `.env` in `reputap-fe/` (or via Nuxt runtime config):

```env
NUXT_PUBLIC_API_BASE=http://localhost:8081
```

## Setup & Development

### 1. Install Dependencies
```bash
bun install
```

### 2. Start Dev Server
```bash
bun run dev
```
Open `http://localhost:3000` in browser.

### 3. Production Build
```bash
bun run build
```

Preview production build:
```bash
bun run preview
```

## Features & Pages

- **Authentication**: JWT login, session persistence, automatic token renewal, route guards.
- **Dashboard**: High-level overview of orders, scan totals, stores, and recent audit activity.
- **Customers**: Manage customer directory, businesses, and contacts.
- **Stores**: Store directory, review URL configurations, and verification.
- **Products**: NFC hardware products catalogue (materials, chip type, pricing).
- **Orders**: Orders listing, item management, and status workflows.
- **Designs**: Order design asset uploads, versioning, and approval flows.
- **Production**: Job tracking, technician assignment, and manufacturing checklists.
- **Quality Control**: Inspection results, checklist verification, and photo evidence.
- **Shipping**: Shipment packing, courier tracking numbers, and delivery status.
- **QR Codes**: Active codes, linked stores, and instance statuses.
- **Analytics**: Real-time scan metrics, daily trends, device breakdown, and store rankings.
- **Audit Logs**: Immutable log of administrative actions and changes.

## Default Admin Credentials
- **Email**: `admin@reputap.id`
- **Password**: `admin123`
