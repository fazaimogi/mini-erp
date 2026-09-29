package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/repositories"
	"github.com/fazasuny/erp-system/internal/routes"
	"github.com/fazasuny/erp-system/internal/services"
)

// Compile-time proof that the GORM implementations satisfy the repository
// contracts used by the services.
var (
	_ repositories.UserRepository    = (*repositories.GormUserRepository)(nil)
	_ repositories.ProductRepository = (*repositories.GormProductRepository)(nil)
)

func newTestRouter() (*gin.Engine, *services.TokenManager) {
	gin.SetMode(gin.TestMode)

	tokens := services.NewTokenManager("smoke-secret", time.Hour)

	authHandler := handlers.NewAuthHandler(
		services.NewAuthService(repositories.NewInMemoryUserRepository(), tokens),
	)
	productRepo := repositories.NewInMemoryProductRepository()
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	inventoryHandler := handlers.NewInventoryHandler(
		services.NewInventoryService(productService, repositories.NewInMemoryInventoryRepository(productRepo)),
	)

	router := gin.New()
	api := router.Group("/api")
	routes.RegisterAuthRoutes(api, authHandler, tokens)
	routes.RegisterProductRoutes(api, productHandler, tokens)
	routes.RegisterInventoryRoutes(api, inventoryHandler, tokens)

	return router, tokens
}

func do(t *testing.T, router *gin.Engine, method, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func decode(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var payload map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("response is not JSON: %v (%s)", err, recorder.Body.String())
	}
	return payload
}

func errorCode(payload map[string]interface{}) string {
	errObj, ok := payload["error"].(map[string]interface{})
	if !ok {
		return ""
	}
	code, _ := errObj["code"].(string)
	return code
}

func TestAuthFlowEndToEnd(t *testing.T) {
	router, _ := newTestRouter()

	register := do(t, router, http.MethodPost, "/api/auth/register", "",
		`{"name":"Budi","email":"budi@example.com","password":"secret12345"}`)
	if register.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d (%s)", register.Code, register.Body.String())
	}

	registered := decode(t, register)
	if data := registered["data"].(map[string]interface{}); data["role"] != "staff" {
		t.Errorf("register: expected staff role, got %v", data["role"])
	}
	if strings.Contains(register.Body.String(), "password") {
		t.Errorf("register: response leaks password field: %s", register.Body.String())
	}

	duplicate := do(t, router, http.MethodPost, "/api/auth/register", "",
		`{"name":"Other","email":"budi@example.com","password":"secret12345"}`)
	if duplicate.Code != http.StatusConflict || errorCode(decode(t, duplicate)) != "EMAIL_ALREADY_EXISTS" {
		t.Errorf("duplicate register: expected 409 EMAIL_ALREADY_EXISTS, got %d (%s)",
			duplicate.Code, duplicate.Body.String())
	}

	badLogin := do(t, router, http.MethodPost, "/api/auth/login", "",
		`{"email":"budi@example.com","password":"wrong-password"}`)
	if badLogin.Code != http.StatusUnauthorized || errorCode(decode(t, badLogin)) != "INVALID_CREDENTIALS" {
		t.Errorf("bad login: expected 401 INVALID_CREDENTIALS, got %d (%s)",
			badLogin.Code, badLogin.Body.String())
	}

	login := do(t, router, http.MethodPost, "/api/auth/login", "",
		`{"email":"budi@example.com","password":"secret12345"}`)
	if login.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d (%s)", login.Code, login.Body.String())
	}

	loginData := decode(t, login)["data"].(map[string]interface{})
	token, ok := loginData["token"].(string)
	if !ok || token == "" {
		t.Fatalf("login: no token in response: %s", login.Body.String())
	}

	anonymousMe := do(t, router, http.MethodGet, "/api/auth/me", "", "")
	if anonymousMe.Code != http.StatusUnauthorized || errorCode(decode(t, anonymousMe)) != "UNAUTHORIZED" {
		t.Errorf("anonymous me: expected 401 UNAUTHORIZED, got %d (%s)",
			anonymousMe.Code, anonymousMe.Body.String())
	}

	me := do(t, router, http.MethodGet, "/api/auth/me", token, "")
	if me.Code != http.StatusOK {
		t.Fatalf("me: expected 200, got %d (%s)", me.Code, me.Body.String())
	}
	if data := decode(t, me)["data"].(map[string]interface{}); data["email"] != "budi@example.com" {
		t.Errorf("me: unexpected payload %v", data)
	}

	logout := do(t, router, http.MethodPost, "/api/auth/logout", token, "")
	if logout.Code != http.StatusOK {
		t.Errorf("logout: expected 200, got %d (%s)", logout.Code, logout.Body.String())
	}

	garbage := do(t, router, http.MethodGet, "/api/auth/me", "not-a-jwt", "")
	if garbage.Code != http.StatusUnauthorized || errorCode(decode(t, garbage)) != "INVALID_TOKEN" {
		t.Errorf("tampered token: expected 401 INVALID_TOKEN, got %d (%s)",
			garbage.Code, garbage.Body.String())
	}
}

func TestProductRoutesRequireAuthAndAdminRole(t *testing.T) {
	router, tokens := newTestRouter()

	staffToken, err := tokens.Generate(2, "staff@example.com", "staff")
	if err != nil {
		t.Fatalf("token generation failed: %v", err)
	}
	adminToken, err := tokens.Generate(1, "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("token generation failed: %v", err)
	}

	anonymous := do(t, router, http.MethodGet, "/api/products", "", "")
	if anonymous.Code != http.StatusUnauthorized {
		t.Errorf("anonymous list: expected 401, got %d (%s)", anonymous.Code, anonymous.Body.String())
	}

	list := do(t, router, http.MethodGet, "/api/products", staffToken, "")
	if list.Code != http.StatusOK {
		t.Errorf("staff list: expected 200, got %d (%s)", list.Code, list.Body.String())
	}

	payload := `{"name":"Keyboard","sku":"KEY-001","category_id":1,"price":250000,"minimum_stock":5}`

	staffCreate := do(t, router, http.MethodPost, "/api/products", staffToken, payload)
	if staffCreate.Code != http.StatusForbidden || errorCode(decode(t, staffCreate)) != "FORBIDDEN" {
		t.Errorf("staff create: expected 403 FORBIDDEN, got %d (%s)",
			staffCreate.Code, staffCreate.Body.String())
	}

	adminCreate := do(t, router, http.MethodPost, "/api/products", adminToken, payload)
	if adminCreate.Code != http.StatusCreated {
		t.Errorf("admin create: expected 201, got %d (%s)", adminCreate.Code, adminCreate.Body.String())
	}
}

func TestInventoryFlowEndToEnd(t *testing.T) {
	router, tokens := newTestRouter()

	staffToken, err := tokens.Generate(2, "staff@example.com", "staff")
	if err != nil {
		t.Fatalf("token generation failed: %v", err)
	}

	// Product 1 is seeded with stock 10, minimum stock 3.
	list := do(t, router, http.MethodGet, "/api/inventory", staffToken, "")
	if list.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d (%s)", list.Code, list.Body.String())
	}
	meta := decode(t, list)["meta"].(map[string]interface{})
	if meta["total"].(float64) != 3 {
		t.Errorf("expected 3 inventory rows, got %v", meta["total"])
	}

	anonymous := do(t, router, http.MethodPost, "/api/inventory/stock-in", "", `{"product_id":1,"quantity":5}`)
	if anonymous.Code != http.StatusUnauthorized {
		t.Errorf("anonymous stock in: expected 401, got %d (%s)", anonymous.Code, anonymous.Body.String())
	}

	// PRD section 9.1: 10 + 5 = 15.
	stockIn := do(t, router, http.MethodPost, "/api/inventory/stock-in", staffToken, `{"product_id":1,"quantity":5,"reference":"PO-001"}`)
	if stockIn.Code != http.StatusOK {
		t.Fatalf("stock in: expected 200, got %d (%s)", stockIn.Code, stockIn.Body.String())
	}
	data := decode(t, stockIn)["data"].(map[string]interface{})
	if data["type"] != "IN" || data["previous_stock"].(float64) != 10 || data["current_stock"].(float64) != 15 {
		t.Errorf("unexpected stock in payload: %v", data)
	}

	// PRD section 9.2: 15 - 3 = 12.
	stockOut := do(t, router, http.MethodPost, "/api/inventory/stock-out", staffToken, `{"product_id":1,"quantity":3,"reference":"INV-100"}`)
	if stockOut.Code != http.StatusOK {
		t.Fatalf("stock out: expected 200, got %d (%s)", stockOut.Code, stockOut.Body.String())
	}
	data = decode(t, stockOut)["data"].(map[string]interface{})
	if data["previous_stock"].(float64) != 15 || data["current_stock"].(float64) != 12 {
		t.Errorf("unexpected stock out payload: %v", data)
	}

	// PRD section 9.3: requesting more than available must be refused.
	excessive := do(t, router, http.MethodPost, "/api/inventory/stock-out", staffToken, `{"product_id":1,"quantity":13}`)
	if excessive.Code != http.StatusConflict || errorCode(decode(t, excessive)) != "INSUFFICIENT_STOCK" {
		t.Errorf("excessive stock out: expected 409 INSUFFICIENT_STOCK, got %d (%s)",
			excessive.Code, excessive.Body.String())
	}

	zero := do(t, router, http.MethodPost, "/api/inventory/stock-in", staffToken, `{"product_id":1,"quantity":0}`)
	if zero.Code != http.StatusUnprocessableEntity || errorCode(decode(t, zero)) != "INVALID_QUANTITY" {
		t.Errorf("zero quantity: expected 422 INVALID_QUANTITY, got %d (%s)", zero.Code, zero.Body.String())
	}

	unknown := do(t, router, http.MethodPost, "/api/inventory/stock-in", staffToken, `{"product_id":9999,"quantity":1}`)
	if unknown.Code != http.StatusNotFound || errorCode(decode(t, unknown)) != "PRODUCT_NOT_FOUND" {
		t.Errorf("unknown product: expected 404 PRODUCT_NOT_FOUND, got %d (%s)", unknown.Code, unknown.Body.String())
	}

	history := do(t, router, http.MethodGet, "/api/inventory/history?product_id=1", staffToken, "")
	if history.Code != http.StatusOK {
		t.Fatalf("history: expected 200, got %d (%s)", history.Code, history.Body.String())
	}
	payload := decode(t, history)
	rows := payload["data"].([]interface{})
	if len(rows) != 2 {
		t.Fatalf("expected 2 history rows (rejected movements are not recorded), got %d", len(rows))
	}
	newest := rows[0].(map[string]interface{})
	if newest["type"] != "OUT" || newest["previous_stock"].(float64) != 15 || newest["current_stock"].(float64) != 12 {
		t.Errorf("expected newest row to be the OUT 15->12, got %v", newest)
	}
	if newest["created_by"].(float64) != 2 {
		t.Errorf("expected created_by to record the caller, got %v", newest["created_by"])
	}

	detail := do(t, router, http.MethodGet, "/api/inventory/1", staffToken, "")
	if detail.Code != http.StatusOK {
		t.Fatalf("detail: expected 200, got %d (%s)", detail.Code, detail.Body.String())
	}
	if stock := decode(t, detail)["data"].(map[string]interface{})["stock"].(float64); stock != 12 {
		t.Errorf("expected detail stock 12, got %v", stock)
	}

	lowStock := do(t, router, http.MethodGet, "/api/products/low-stock", staffToken, "")
	if lowStock.Code != http.StatusOK {
		t.Fatalf("low stock: expected 200, got %d (%s)", lowStock.Code, lowStock.Body.String())
	}
	rows = decode(t, lowStock)["data"].([]interface{})
	if len(rows) != 1 {
		t.Fatalf("expected only the seat product to be low stock, got %d rows", len(rows))
	}
	if rows[0].(map[string]interface{})["sku"] != "KURSI-001" {
		t.Errorf("expected KURSI-001 to be low stock, got %v", rows[0])
	}
}
