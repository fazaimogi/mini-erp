# Testing Documentation

Complete test suite: 41 tests (27 backend + 14 frontend), all passing.

## Backend Testing

### Running tests

```bash
cd backend
go test -v ./...          # All tests
go test -v ./internal/... # Only internal packages
go test -run TestAuth     # Specific test pattern
```

### Test packages

**`internal/routes/routes_test.go`** (3 tests)
- End-to-end integration tests with real database
- Tests: Auth flow, RBAC enforcement, inventory transactions

**`internal/services/auth_service_test.go`** (8 tests)
- Registration: password hashing, default role assignment
- Registration validation: email format, password length, name length
- Login: token generation, invalid credentials
- Token: expiry validation, tampering detection
- User retrieval with role

**`internal/services/product_service_test.go`** (10 tests)
- Create: assignment of ID, timestamps
- Duplicate SKU rejection
- Input validation: name, SKU, category, price (negative), stock
- Get by ID: correct retrieval
- Pagination: offset/limit, total count
- Update: field changes, SKU uniqueness across products
- Delete: actual removal

**`internal/services/inventory_service_test.go`** (6 tests)
- Stock in: adds quantity, records history
- Stock out: deducts quantity, records history
- Insufficient stock rejection
- Exact available deduction allowed
- Non-positive quantity rejection
- Movement on unknown product returns 404
- History filtering by product + pagination
- Inventory list/detail expose stock

### Coverage areas

- **Auth:** Registration validation, bcrypt hashing, JWT signing/parsing, role assignment
- **Products:** CRUD operations, SKU uniqueness, business rules (negative prices)
- **Inventory:** Stock movements, transaction safety via `SELECT ... FOR UPDATE`
- **Routes:** RBAC middleware, auth middleware, error responses
- **Database:** Transactions, migrations, constraints

### Key fixtures

```go
// Inventory tests set up:
// Product 1: "Laptop Asus XYZ" - stock 10, min 3
// Product 3: "Kursi Kantor Ergonomis" - stock 2 (low)
// Tests verify movements against these fixed products
```

### Running specific tests

```bash
# Auth tests only
go test -v -run TestAuth ./internal/services/

# Inventory out movements
go test -v -run TestStockOut ./internal/services/

# Verbose output with timing
go test -v -count=1 ./...  # -count=1 disables caching
```

## Frontend Testing

### Running tests

```bash
cd frontend
npm test              # Run all tests
npm run test:watch   # Watch mode
npm test -- --coverage  # Coverage report (when setup)
```

### Test files

**`__tests__/api.test.ts`** (3 tests)
- Authorization header with token
- Error handling (non-ok response)
- Network error handling

**`__tests__/Header.test.tsx`** (5 tests)
- Title rendering
- User name display (regex match for DOM structure)
- Null user handling
- Action slot rendering
- Role object support

**`__tests__/Sidebar.test.tsx`** (6 tests)
- Navigation links present
- Admin-only links: Categories, Users (shown for admin, hidden for staff)
- Logout button
- Null user handling
- Role as string or object

### Coverage areas

- **API client:** JWT auth header, error responses, network errors
- **Components:** Conditional rendering, role-based visibility
- **RBAC:** Admin/Staff link visibility in UI
- **User handling:** null users, role formats

## Writing new tests

### Backend (Go)

Pattern:

```go
func TestFeatureName(t *testing.T) {
    // Setup
    service := setupService()
    input := validInput()
    
    // Act
    result, err := service.Create(input)
    
    // Assert
    if err != nil {
        t.Fatalf("expected success, got error: %v", err)
    }
    if result.ID == 0 {
        t.Error("expected assigned ID")
    }
}
```

**Fixtures:**
- Use in-memory repositories for unit tests
- Use real database for integration tests (in `routes_test.go`)

**Naming:**
- `TestFeatureName` for happy path
- `TestFeatureName_ErrorCase` for validation/edge cases
- `TestFeatureName/subcase` for table-driven tests

### Frontend (Next.js + React Testing Library)

Pattern:

```typescript
describe('ComponentName', () => {
  it('renders correctly', () => {
    render(<Component prop="value" />);
    expect(screen.getByText('Expected text')).toBeInTheDocument();
  });

  it('handles user interaction', async () => {
    const { user } = render(<Component />);
    await user.click(screen.getByRole('button', { name: /click/i }));
    expect(screen.getByText('After click')).toBeInTheDocument();
  });
});
```

**Queries (in order of preference):**
1. `getByRole` (accessibility) — `screen.getByRole('button', { name: /Login/i })`
2. `getByLabelText` (forms) — `screen.getByLabelText('Email')`
3. `getByText` (content) — `screen.getByText(/exact text/)`
4. `getByPlaceholderText` (inputs) — `screen.getByPlaceholderText('Enter name')`
5. Regex for partial matches: `/text/i`

**Mocking:**
```typescript
jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
}));

global.fetch = jest.fn(() =>
  Promise.resolve({ ok: true, json: () => Promise.resolve({...}) })
);
```

## CI/CD integration

### GitHub Actions (example)

```yaml
name: Tests
on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: erp_system_test
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: cd backend && go test -v ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: cd frontend && npm ci && npm test
```

## Debugging tests

### Backend

```bash
# Verbose output
go test -v -run TestName ./...

# Verbose + trace
go test -v -trace=trace.out ./...
go tool trace trace.out

# Debug failing test
go test -run TestName -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Frontend

```bash
# Run in band (no parallel execution)
npm test -- --runInBand

# Only failed tests
npm test -- --onlyChanged

# Update snapshots (if using snapshot testing)
npm test -- -u

# Debug in browser
node --inspect-brk node_modules/.bin/jest --runInBand
# Then open chrome://inspect
```

## Performance benchmarks

### Backend

Create `*_bench_test.go`:

```go
func BenchmarkCreateProduct(b *testing.B) {
    service := setupService()
    input := validInput()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.Create(input)
    }
}
```

Run:
```bash
go test -bench=. -benchmem ./internal/services
```

### Frontend

Use React DevTools Profiler:
1. Open Chrome DevTools → React tab
2. Click "Profiler" tab
3. Record interactions
4. Check render times

## Known issues & gotchas

### Backend

- Tests use in-memory repos; integration tests need real PostgreSQL
- `SELECT ... FOR UPDATE` requires real DB transactions
- Migrations run on server startup (tests must reset DB)

### Frontend

- Component rendering tests check final DOM (not intermediate states)
- Regex matchers needed for multipart text (e.g., "Welcome, John Doe")
- Mocking `useRouter` from Next.js is required for page tests

## Test maintenance

### Add tests for:
- New feature endpoints
- Validation rules (always add negative tests)
- Error cases (404, 400, 409, etc.)
- RBAC changes (admin vs staff visibility)
- Business logic (stock deduction, inventory history)

### Remove tests for:
- Deleted features
- Forwarding/wiring (if both endpoints tested, don't test wrapper)
- Tautologies (e.g., testing that a mock returns its preset value)

### Review tests when:
- Refactoring services (update test setup)
- Changing API responses (update assertions)
- Adding new roles/permissions (add RBAC tests)

## Coverage targets

- **Backend:** 70%+ line coverage (focus on services)
- **Frontend:** 50%+ (components and API client)
- **Critical paths:** 100% (auth, stock deduction, RBAC)

Run coverage report:

```bash
# Backend
go test -cover ./...
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Frontend (requires setup)
npm test -- --coverage
```
