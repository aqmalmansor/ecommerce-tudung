package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"be/config"
	"be/enums"
	"be/handlers"
	"be/middleware"
	"be/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductOption{},
		&models.ProductOptionValue{},
		&models.ProductVariant{},
		&models.ProductImage{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ph := handlers.NewProductHandler(db)

	products := r.Group("/api/products")
	products.GET("", ph.ListProducts)
	products.GET("/:id", ph.GetProduct)

	admin := products.Group("")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireAdmin(db))
	admin.POST("", ph.CreateProduct)
	admin.PUT("/:id", ph.UpdateProduct)
	admin.DELETE("/:id", ph.DeleteProduct)
	admin.POST("/:id/variants", ph.AddVariant)
	admin.PUT("/:id/variants/:variantId", ph.UpdateVariant)
	admin.DELETE("/:id/variants/:variantId", ph.DeleteVariant)

	return r
}

func seedUser(t *testing.T, db *gorm.DB, email, role string) models.User {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.MinCost)
	u := models.User{Email: email, Password: string(hash), Name: role, Role: role}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}

func bearerToken(t *testing.T, userID uint, email string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"email":   email,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return "Bearer " + tok
}

func seedProduct(t *testing.T, db *gorm.DB, handle, status string) models.Product {
	t.Helper()
	p := models.Product{Title: handle, Handle: handle, Status: status}
	if err := db.Create(&p).Error; err != nil {
		t.Fatalf("seed product: %v", err)
	}
	return p
}

func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func do(r *gin.Engine, method, path string, body *bytes.Reader, token string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- List ---

func TestListProducts_DefaultsToActive(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	seedProduct(t, db, "active-product", enums.ProductStatusActive.String())
	seedProduct(t, db, "draft-product", enums.ProductStatusDraft.String())

	w := do(r, http.MethodGet, "/api/products", nil, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if got := len(resp["products"].([]any)); got != 1 {
		t.Errorf("expected 1 active product, got %d", got)
	}
}

func TestListProducts_FilterByStatus(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	seedProduct(t, db, "active-p", enums.ProductStatusActive.String())
	seedProduct(t, db, "draft-p", enums.ProductStatusDraft.String())

	w := do(r, http.MethodGet, "/api/products?status="+enums.ProductStatusDraft.String(), nil, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if got := len(resp["products"].([]any)); got != 1 {
		t.Errorf("expected 1 draft product, got %d", got)
	}
}

// --- Get ---

func TestGetProduct_Found(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	p := seedProduct(t, db, "my-product", enums.ProductStatusActive.String())

	w := do(r, http.MethodGet, fmt.Sprintf("/api/products/%d", p.ID), nil, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := do(r, http.MethodGet, "/api/products/99999", nil, "")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- Create ---

func TestCreateProduct_Admin_Success(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())

	payload := map[string]any{
		"title":  "Tudung Bawal",
		"handle": "tudung-bawal",
		"status": enums.ProductStatusActive.String(),
		"options": []map[string]any{
			{"name": "Color", "values": []string{"Red", "Blue"}},
			{"name": "Size", "values": []string{"S", "M", "L"}},
		},
		"variants": []map[string]any{
			{"sku": "TB-RED-M", "price": 29.90, "inventory_quantity": 50, "option1": "Red", "option2": "M"},
		},
	}
	w := do(r, http.MethodPost, "/api/products", jsonBody(payload), bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	p := resp["product"].(map[string]any)
	if p["title"] != "Tudung Bawal" {
		t.Errorf("unexpected title: %v", p["title"])
	}
	if p["handle"] != "tudung-bawal" {
		t.Errorf("unexpected handle: %v", p["handle"])
	}
}

func TestCreateProduct_Unauthorized(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := do(r, http.MethodPost, "/api/products", jsonBody(map[string]any{"title": "X", "handle": "x"}), "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateProduct_Forbidden_Client(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	client := seedUser(t, db, "client@test.com", enums.Client.String())

	w := do(r, http.MethodPost, "/api/products", jsonBody(map[string]any{"title": "X", "handle": "x"}), bearerToken(t, client.ID, client.Email))

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestCreateProduct_InvalidStatus(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())

	w := do(r, http.MethodPost, "/api/products",
		jsonBody(map[string]any{"title": "X", "handle": "x", "status": "invalid"}),
		bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Update ---

func TestUpdateProduct_Status(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())
	p := seedProduct(t, db, "my-product", enums.ProductStatusDraft.String())

	w := do(r, http.MethodPut, fmt.Sprintf("/api/products/%d", p.ID),
		jsonBody(map[string]any{"status": enums.ProductStatusActive.String()}),
		bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if got := resp["product"].(map[string]any)["status"]; got != enums.ProductStatusActive.String() {
		t.Errorf("expected status %s, got %v", enums.ProductStatusActive.String(), got)
	}
}

// --- Delete ---

func TestDeleteProduct(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())
	p := seedProduct(t, db, "my-product", enums.ProductStatusActive.String())

	w := do(r, http.MethodDelete, fmt.Sprintf("/api/products/%d", p.ID), nil, bearerToken(t, admin.ID, admin.Email))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	w2 := do(r, http.MethodGet, fmt.Sprintf("/api/products/%d", p.ID), nil, "")
	if w2.Code != http.StatusNotFound {
		t.Errorf("expected 404 after soft-delete, got %d", w2.Code)
	}
}

// --- Variants ---

func TestAddVariant_Success(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())
	p := seedProduct(t, db, "my-product", enums.ProductStatusActive.String())

	payload := map[string]any{"sku": "TB-BLUE-L", "price": 35.00, "inventory_quantity": 20, "option1": "Blue", "option2": "L"}
	w := do(r, http.MethodPost, fmt.Sprintf("/api/products/%d/variants", p.ID), jsonBody(payload), bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAddVariant_DuplicateSKU(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())
	p := seedProduct(t, db, "my-product", enums.ProductStatusActive.String())
	db.Create(&models.ProductVariant{ProductID: p.ID, SKU: "DUPE-SKU", Price: 10})

	w := do(r, http.MethodPost, fmt.Sprintf("/api/products/%d/variants", p.ID),
		jsonBody(map[string]any{"sku": "DUPE-SKU", "price": 20.00}),
		bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 on duplicate SKU, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateVariant_InventoryQty(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())
	p := seedProduct(t, db, "my-product", enums.ProductStatusActive.String())
	v := models.ProductVariant{ProductID: p.ID, SKU: "TB-RED-S", Price: 29.90, InventoryQty: 10}
	db.Create(&v)

	w := do(r, http.MethodPut, fmt.Sprintf("/api/products/%d/variants/%d", p.ID, v.ID),
		jsonBody(map[string]any{"inventory_quantity": 99}),
		bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var updated models.ProductVariant
	db.First(&updated, v.ID)
	if updated.InventoryQty != 99 {
		t.Errorf("expected inventory 99, got %d", updated.InventoryQty)
	}
}

func TestDeleteVariant(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())
	p := seedProduct(t, db, "my-product", enums.ProductStatusActive.String())
	v := models.ProductVariant{ProductID: p.ID, SKU: "TB-DEL", Price: 10.00}
	db.Create(&v)

	w := do(r, http.MethodDelete, fmt.Sprintf("/api/products/%d/variants/%d", p.ID, v.ID), nil, bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateVariant_WrongProduct(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)
	admin := seedUser(t, db, "admin@test.com", enums.Admin.String())
	p1 := seedProduct(t, db, "product-1", enums.ProductStatusActive.String())
	p2 := seedProduct(t, db, "product-2", enums.ProductStatusActive.String())
	v := models.ProductVariant{ProductID: p1.ID, SKU: "P1-SKU", Price: 10.00}
	db.Create(&v)

	w := do(r, http.MethodPut, fmt.Sprintf("/api/products/%d/variants/%d", p2.ID, v.ID),
		jsonBody(map[string]any{"price": 99.0}),
		bearerToken(t, admin.ID, admin.Email))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
