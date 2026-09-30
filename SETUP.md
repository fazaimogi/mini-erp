# Setup Guide

Step-by-step installation and configuration.

## Prerequisites

- **Backend:** Go 1.21+, PostgreSQL 12+
- **Frontend:** Node.js 18+, npm 9+
- **Optional:** Docker + Docker Compose

## Backend Setup

### 1. Environment

```bash
cd backend
cp .env.example .env
```

Edit `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=erp_system
DB_SSLMODE=disable
JWT_SECRET=your_random_secret_key_32_chars_minimum
JWT_EXPIRY_HOURS=24
```

### 2. PostgreSQL

**Option A: Local installation**

```bash
# macOS
brew install postgresql@15
brew services start postgresql@15

# Linux (Ubuntu/Debian)
sudo apt-get install postgresql postgresql-contrib
sudo service postgresql start

# Windows
# Download from https://www.postgresql.org/download/windows/
```

**Option B: Docker**

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_secure_password \
  -e POSTGRES_DB=erp_system \
  -p 5432:5432 \
  postgres:15

# Or use docker-compose:
docker-compose up -d postgres
```

### 3. Create database

```bash
createdb -U postgres erp_system
# or
psql -U postgres -c "CREATE DATABASE erp_system;"
```

### 4. Install Go dependencies

```bash
cd backend
go mod download
```

### 5. Run migrations

Migrations run automatically on server startup. Or manually:

```bash
cd backend
go run ./cmd/server
# Watch for "migration complete" in logs
```

### 6. Start backend

```bash
cd backend
go run ./cmd/server
# Output: listening on :8080
```

**Verify:**

```bash
curl http://localhost:8080
# {"message":"ERP API is running"}
```

### Testing backend

```bash
cd backend
go test -v ./...
# 27 tests should pass
```

## Frontend Setup

### 1. Install dependencies

```bash
cd frontend
npm install
```

### 2. Environment

Create `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### 3. Start development server

```bash
npm run dev
# Output: ▲ Next.js 16.3.6
# Ready in 2.5s on http://localhost:3000
```

Open http://localhost:3000

### 4. Build for production

```bash
npm run build
npm start
# Production build ready on http://localhost:3000
```

### Testing frontend

```bash
npm test
# 14 tests should pass
```

## Docker (all-in-one)

### Start services

```bash
docker-compose up -d
```

Services:
- PostgreSQL: localhost:5432
- Backend: localhost:8080
- Frontend: localhost:3000 (run separately: `npm run dev`)

### View logs

```bash
docker-compose logs -f backend
docker-compose logs -f postgres
```

### Stop services

```bash
docker-compose down
```

## First use

### 1. Register user

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Admin",
    "email": "admin@example.com",
    "password": "SecurePass123!"
  }'
```

Response:
```json
{
  "data": {
    "user": {"id": 1, "name": "Admin", "role": "staff"},
    "token": "eyJ..."
  },
  "message": "Registered successfully"
}
```

### 2. Promote to admin (backend)

```bash
psql -U postgres -d erp_system -c \
  "UPDATE users SET role_id = 1 WHERE email = 'admin@example.com';"
```

### 3. Login frontend

Go to http://localhost:3000/login
- Email: admin@example.com
- Password: SecurePass123!

### 4. Create first product

1. Dashboard → Products → "+ New Product"
2. Fill form:
   - Name: "Laptop"
   - SKU: "LAPTOP-001"
   - Category: (create one first via Categories)
   - Price: 1000000
   - Stock: 10
   - Minimum Stock: 2
3. Save

## Troubleshooting

### Backend won't connect to PostgreSQL

**Error:** `pq: connect: connection refused`

**Fix:**
1. Verify PostgreSQL is running:
   ```bash
   psql -U postgres
   ```
2. Check `.env` DB_* vars match your setup
3. If using Docker:
   ```bash
   docker ps | grep postgres
   docker exec <container_id> pg_isready
   ```

### Frontend can't reach backend

**Error:** `Failed to load data` in browser

**Fix:**
1. Verify backend is running:
   ```bash
   curl http://localhost:8080
   ```
2. Check `NEXT_PUBLIC_API_URL` in `.env.local`
3. Check browser DevTools → Network → API calls for CORS errors
4. Backend CORS is enabled for `http://localhost:3000`

### JWT errors on login

**Error:** `INVALID_CREDENTIALS` for correct email/password

**Fix:**
1. Verify `JWT_SECRET` is set in backend `.env`
2. Restart backend after changing `JWT_SECRET`
3. Clear browser localStorage:
   ```javascript
   localStorage.clear()
   ```

### Tests fail

**Backend tests:**
```bash
cd backend
go test -v ./...  # Should show 27 PASS
```

**Frontend tests:**
```bash
cd frontend
npm test  # Should show 14 PASS
```

If tests fail, check:
1. All dependencies installed
2. No port conflicts (8080, 3000, 5432)
3. Fresh database (`dropdb erp_system && createdb erp_system` for backend tests)

## Performance tuning

### PostgreSQL

For development, defaults are fine. For production:

```sql
-- Connection pooling
-- Use PgBouncer: https://www.pgbouncer.org/

-- Indexes (auto-created by migrations)
-- Check: SELECT * FROM pg_stat_user_indexes;

-- Slow query log
ALTER SYSTEM SET log_min_duration_statement = 1000;  -- 1s threshold
SELECT pg_reload_conf();
```

### Backend

- Increase connection pool: `DB_MAX_OPEN_CONNS=25`
- Enable query caching in repositories (future)
- Use Redis for session storage (future)

### Frontend

- Enable SWR caching in API client (future)
- Code splitting via Next.js automatic routes

## Monitoring

### Backend health check

```bash
curl http://localhost:8080  # Should return 200 OK
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN"  # Verify auth
```

### Database health

```bash
psql -U postgres -d erp_system -c \
  "SELECT COUNT(*) FROM users; SELECT COUNT(*) FROM products;"
```

### Frontend health

Open http://localhost:3000 in browser; check:
- Login page loads
- Can login with valid credentials
- Dashboard displays

## Cleanup

### Stop everything

```bash
# Backend: Ctrl+C in terminal
# Frontend: Ctrl+C in terminal
# Docker: docker-compose down
```

### Reset database

```bash
dropdb -U postgres erp_system
createdb -U postgres erp_system
# Restart backend (runs migrations)
```

### Uninstall

```bash
# Backend: rm -rf backend
# Frontend: rm -rf frontend
# PostgreSQL: brew uninstall postgresql@15  (or apt-get remove)
```

## Support

- **Issues:** Check existing GitHub issues
- **Logs:** `backend/logs/` and browser DevTools Console
- **PRD:** Reference `PRD.txt` for feature requirements
