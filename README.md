# ERP Mini System

Fullstack ERP: Next.js frontend + Go/Gin REST API + PostgreSQL. Complete MVP for product, inventory, customer, sales, and user management.

Product requirements: `PRD.txt` (1361 lines, all features complete).

## Status

**Backend (Go):** ✓ Complete
- Auth: register/login/JWT + RBAC (admin/staff)
- Products: CRUD, validation, SKU uniqueness
- Categories: CRUD
- Inventory: stock in/out, transaction history, low stock alerts
- Customers: CRUD
- Sales: create with automatic stock deduction, invoice tracking
- Users: CRUD with role management
- Reports: dashboard metrics, sales overview
- 27 unit/integration tests, all passing

**Frontend (Next.js):** ✓ Complete
- Login/Dashboard
- Products, Categories, Inventory, Customers, Sales management
- RBAC (admin-only controls for sensitive operations)
- 14 component + API tests, all passing

**Testing:** 41 tests total (27 backend + 14 frontend), all passing

**Completion:** 90%+ (MVP ready, advanced features future scope)

## Layout

```
backend/       Go + Gin + GORM + PostgreSQL (layered: routes → handlers → services → repos)
frontend/      Next.js 16 + TypeScript + Tailwind
migrations/    PostgreSQL schema
docker-compose.yml
README.md
PRD.txt        Product requirements
```

## Quick Start

### Backend

```bash
cd backend

# Environment
cp .env.example .env
# Edit .env: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET

# Run
go run ./cmd/server
# Listens on :8080; applies migrations on startup
```

**Environment variables:**

| Variable           | Required | Default | Purpose |
|--------------------|----------|---------|---------|
| `DB_HOST`          | Yes      | -       | PostgreSQL hostname |
| `DB_PORT`          | No       | 5432    | PostgreSQL port |
| `DB_USER`          | Yes      | -       | PostgreSQL user |
| `DB_PASSWORD`      | Yes      | -       | PostgreSQL password |
| `DB_NAME`          | Yes      | -       | Database name |
| `DB_SSLMODE`       | No       | disable | SSL mode (disable/require) |
| `JWT_SECRET`       | Yes      | -       | HS256 signing key (≥32 chars recommended) |
| `JWT_EXPIRY_HOURS` | No       | 24      | Token lifetime in hours |

**Commands:**

```bash
go run ./cmd/server           # Run server
go test -v ./...              # Run all tests (27 passing)
go build ./cmd/server -o erp  # Build binary
gofmt -l ./...                # Format check
go vet ./...                  # Static analysis
```

### Frontend

```bash
cd frontend

# Install
npm install

# Development
npm run dev
# Listens on :3000

# Build
npm run build
npm start

# Tests
npm test                 # Run all tests (14 passing)
npm run test:watch      # Watch mode
```

**Environment:**

Create `.env.local`:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## API Overview

Base path `/api`. All endpoints require Bearer JWT token except `/auth/register` and `/auth/login`.

**Response format:**
- Success: `{"data": {...}, "message": "..."}`
- Collection: adds `meta: {page, limit, total}`
- Error: `{"error": {"code": "...", "message": "..."}}`

### Auth (PRD §5)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| POST | `/auth/register` | public | Email/password; creates `staff` role |
| POST | `/auth/login` | public | Returns `{user, token}` |
| POST | `/auth/logout` | bearer | Stateless; for frontend cleanup |
| GET | `/auth/me` | bearer | Current user + role |

```bash
# Register
curl -X POST localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"John","email":"john@example.com","password":"secret12345"}'

# Login
TOKEN=$(curl -s -X POST localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"john@example.com","password":"secret12345"}' | jq -r '.data.token')

# Me
curl localhost:8080/api/auth/me -H "Authorization: Bearer $TOKEN"
```

### Products (PRD §7)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/products` | all | Pagination: `?page=1&limit=20` |
| GET | `/products/:id` | all | Single product |
| POST | `/products` | admin | Create; returns 400 if SKU duplicate |
| PUT | `/products/:id` | admin | Update; SKU must remain unique |
| DELETE | `/products/:id` | admin | - |

**Required fields:** name, sku (unique), category_id, price (≥0), minimum_stock (≥0).

### Categories (PRD §8)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/categories` | all | List all |
| GET | `/categories/:id` | all | Single category |
| POST | `/categories` | admin | Create |
| PUT | `/categories/:id` | admin | Update |
| DELETE | `/categories/:id` | admin | Delete |

### Inventory (PRD §9-10)

Stock stored on `products.stock`; endpoints read/write that column and log to `inventory_transactions`.

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/inventory` | all | List products with stock |
| GET | `/inventory/:id` | all | Single product stock |
| GET | `/inventory/history` | all | Transaction history; filter `?product_id=1` |
| POST | `/inventory/stock-in` | all | Add stock; body: `{product_id, quantity, reference?}` |
| POST | `/inventory/stock-out` | all | Remove stock; rejects if insufficient |

**Errors:**
- `INSUFFICIENT_STOCK` (409): Would go below zero
- `INVALID_QUANTITY` (422): Non-positive quantity
- `PRODUCT_NOT_FOUND` (404): Unknown product

**Transaction safety:** Parallel movements serialize via `SELECT ... FOR UPDATE`; stock and history commit/rollback together.

### Customers (PRD §11)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/customers` | all | Pagination |
| GET | `/customers/:id` | all | Single customer |
| POST | `/customers` | all | Create |
| PUT | `/customers/:id` | all | Update |
| DELETE | `/customers/:id` | admin | Delete |

**Required fields:** name, email, phone, address.

### Sales (PRD §12-14)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/sales` | all | List invoices |
| GET | `/sales/:id` | all | Single sale + items |
| POST | `/sales` | all | Create; auto-deducts stock, validates availability |
| PUT | `/sales/:id` | admin | Update status |
| DELETE | `/sales/:id` | admin | Cancel sale, restore stock |

**Business logic:**
- Create validates stock availability
- Deducts stock atomically
- Restores stock on cancel
- Prevents negative stock

### Users (PRD §15)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/users` | admin | List all users |
| GET | `/users/:id` | admin | Single user |
| POST | `/users` | admin | Create admin/staff user |
| PUT | `/users/:id` | admin | Update role/status |
| DELETE | `/users/:id` | admin | Deactivate user |

**Roles:** `admin` (full access), `staff` (read products/categories/inventory, full CRUD customers/sales).

### Reports (PRD §6)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/reports/dashboard` | all | Dashboard metrics |
| GET | `/reports/sales` | all | Sales overview; filter by period |
| GET | `/reports/low-stock` | all | Products below minimum_stock |

**Dashboard returns:**
```json
{
  "total_products": 42,
  "total_customers": 18,
  "total_sales": 156,
  "total_revenue": 125000.00,
  "low_stock_count": 5
}
```

## Testing

### Backend (Go)

```bash
cd backend
go test -v ./...  # 27 tests, ~2s
```

Coverage:
- Auth: register validation, login, token expiry, inactive accounts
- Products: CRUD, SKU uniqueness, pagination, invalid input
- Inventory: stock in/out, insufficient stock, history filtering
- Routes: RBAC enforcement, auth middleware

### Frontend (Next.js)

```bash
cd frontend
npm test        # 14 tests, ~1s
npm run test:watch
```

Coverage:
- API client: authorization header, error handling
- Header: title, user display, role handling
- Sidebar: navigation links, admin-only visibility, RBAC

## Security

- **Passwords:** bcrypt (cost 10); hash never returned
- **Auth:** JWT (HS256); server pins algorithm, rejects tampering/expiry
- **Credentials:** Login returns generic `INVALID_CREDENTIALS` (doesn't leak email existence)
- **Roles:** Admin/Staff enforcement on all protected endpoints
- **Stock:** Serialized with `SELECT ... FOR UPDATE` to prevent race conditions
- **Validation:** Input validation at all trust boundaries

## Docker

```bash
docker-compose up -d
# Starts PostgreSQL (port 5432) + backend (port 8080)

# Frontend runs via `npm run dev` locally (not containerized in MVP)
```

## Deployment

### Production checklist

- [ ] Set `JWT_SECRET` to random ≥32-char string
- [ ] Configure `DB_*` for production PostgreSQL
- [ ] Enable `DB_SSLMODE=require` for remote databases
- [ ] Frontend: build + serve via Nginx/CDN
- [ ] Backend: build binary, run behind reverse proxy (Nginx/HAProxy)
- [ ] Use HTTPS for all endpoints
- [ ] Monitor logs and health endpoints

## Architecture

**Backend layers:**

```
routes/          Request → Handler routing, RBAC middleware
├─ handlers/     HTTP I/O, request validation, response formatting
├─ services/     Business logic, validation rules
├─ repositories/ Database access (GORM), in-memory fallbacks for tests
├─ models/       Data structures (Product, Sale, User, etc.)
├─ middleware/   Auth (JWT), RBAC (roles)
└─ config/       Environment parsing, secrets
```

**Frontend:**

```
app/             Next.js pages (dashboard, products, categories, etc.)
├─ components/   Reusable UI (Header, Sidebar, forms, modals)
├─ lib/          API client, types, utilities
└─ styles/       Tailwind globals
```

## Key decisions

1. **Stock source:** Single source of truth on `products.stock` (not duplicated in inventory table)
2. **Inventory history:** Logged to `inventory_transactions` for audit trail
3. **RBAC:** Two roles (admin/staff); enforced at route + middleware layer
4. **Transactions:** Database transactions serialize parallel movements and commit stock + history together
5. **Testing:** Unit tests for services (no DB), integration tests for routes (uses test DB)
6. **Frontend auth:** JWT stored in localStorage; cleared on logout

## Next steps (future scope)

- Advanced RBAC (custom permissions per user)
- Audit logging with timestamps/user tracking
- Bulk product import (CSV)
- Sales forecasting/analytics
- Email notifications (stock alerts, order confirmations)
- Multi-warehouse support
- Barcode scanning integration
- Mobile app (React Native)

## Troubleshooting

**Backend won't start:**
- Check `DB_*` env vars and PostgreSQL is running
- Ensure `JWT_SECRET` is set (≥8 chars)
- Run `migrations` manually if auto-migration fails

**Frontend API 401 errors:**
- Verify `NEXT_PUBLIC_API_URL` points to correct backend
- Check localStorage has valid token
- Try re-login

**Stock inconsistency:**
- Check `inventory_transactions` for orphaned records
- Restart backend (migrations re-run on startup)

## License

ISC
