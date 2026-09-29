# ERP Mini System

Fullstack ERP: Next.js frontend + Go/Gin REST API + PostgreSQL.
Product requirements: `PRD.txt`.

## Layout

```
backend/    Go + Gin API (layered: routes -> handlers -> services -> repositories)
frontend/   Next.js app (not started)
```

## Backend setup

```bash
cd backend
cp .env.example .env          # fill DB_* and JWT_SECRET
go run ./cmd/server           # listens on :8080
```

Environment variables:

| Variable           | Purpose                                        |
| ------------------ | ---------------------------------------------- |
| `DB_HOST`..`DB_NAME`, `DB_SSLMODE` | PostgreSQL connection              |
| `JWT_SECRET`       | HS256 signing key; server refuses to start if empty |
| `JWT_EXPIRY_HOURS` | Token lifetime, default `24`                   |

Apply migrations in order (no migration runner yet):

```bash
for f in migrations/*.sql; do psql "$DATABASE_URL" -f "$f"; done
```

## Commands

```bash
go build ./...    # compile
go test ./...     # unit + route tests (no database required)
gofmt -l .        # formatting check
go vet ./...      # static analysis
```

## API

Base path `/api`. Success `{"data":..,"message":..}`, collection adds
`meta{page,limit,total}`, error `{"error":{"code":..,"message":..}}` (PRD section 24).

### Auth (PRD section 5)

| Method | Path                 | Auth   | Notes                                        |
| ------ | -------------------- | ------ | -------------------------------------------- |
| POST   | `/api/auth/register` | public | Always creates a `staff` account              |
| POST   | `/api/auth/login`    | public | Returns `{user, token}`                       |
| POST   | `/api/auth/logout`   | bearer | Stateless JWT ack for the frontend            |
| GET    | `/api/auth/me`       | bearer | Current user + role                           |

```bash
curl -X POST localhost:8080/api/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"budi@example.com","password":"secret12345"}'
curl localhost:8080/api/auth/me -H "Authorization: Bearer $TOKEN"
```

### Products (PRD section 7)

| Method | Path                       | Auth        |
| ------ | -------------------------- | ----------- |
| GET    | `/api/products`            | admin/staff |
| GET    | `/api/products/:id`        | admin/staff |
| GET    | `/api/products/low-stock`  | admin/staff |
| POST   | `/api/products`            | admin       |
| PUT    | `/api/products/:id`        | admin       |
| DELETE | `/api/products/:id`        | admin       |

### Inventory (PRD sections 9-10, 23)

`GET /api/inventory[/:id]` reports the product row, because stock is stored on
`products.stock` (PRD section 7.1) — there is no separate stock table to keep in
sync. Movements record history in `inventory_transactions`.

| Method | Path                          | Auth        |
| ------ | ----------------------------- | ----------- |
| GET    | `/api/inventory`              | admin/staff |
| GET    | `/api/inventory/:id`          | admin/staff |
| GET    | `/api/inventory/history`      | admin/staff |
| POST   | `/api/inventory/stock-in`     | admin/staff |
| POST   | `/api/inventory/stock-out`    | admin/staff |

The product id is in the body, matching the PRD path shape; `product_id` is
required and `reference` is optional. `product_id` is also accepted as a query
filter on `history`.

```bash
curl -X POST localhost:8080/api/inventory/stock-in -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' -d '{"product_id":1,"quantity":5,"reference":"PO-001"}'
curl "localhost:8080/api/inventory/history?product_id=1" -H "Authorization: Bearer $TOKEN"
```

Errors: `INSUFFICIENT_STOCK` (409) when a movement would push stock below zero,
`INVALID_QUANTITY` (422) for a non-positive quantity, `PRODUCT_NOT_FOUND` (404).

Concurrency: the stock row is read with `SELECT ... FOR UPDATE` inside the
movement transaction, so parallel movements serialise instead of both reading
the same `previous_stock`. Stock and its history row commit or roll back
together (PRD section 26.2).

No adjustment endpoint: PRD section 23 does not list one, so `ADJUSTMENT` exists
in the table's CHECK constraint but is never written.

## Security notes

- Passwords hashed with bcrypt (cost 10); `password_hash` never leaves the API.
- Login returns `INVALID_CREDENTIALS` for both unknown email and wrong password.
- JWT signed HS256; `Parse` pins the algorithm and rejects expired or foreign-signed tokens.
- Roles: `admin`, `staff` (PRD section 4). No admin exists until user management
  (PRD section 15) or a manual `UPDATE users SET role_id=...`.

## Progress against the PRD

Done: Day 1-2 (foundation, PostgreSQL + GORM, Product CRUD), Day 3 (auth:
register, login, bcrypt, JWT, auth middleware, role authorization) and Day 4
(inventory: stock in/out, stock validation, inventory history, low stock).

Not started: categories, customers, sales, users, reports APIs; frontend;
Docker; deployment (PRD sections 8, 11-16, 21-22, 27).

### Stock source

`products.stock` is the single source of truth. The `inventory` table exists per
PRD section 18.1 but is deliberately unused: keeping a second copy of the same
number invites drift. Inventory endpoints therefore read and write
`products.stock` and log every change to `inventory_transactions`.