package main

import (
	"api/orders/configs"
	"api/orders/internal/order"
	"api/orders/internal/product"
	"api/orders/internal/user"
	apidb "api/orders/pkg/db"
	apijwt "api/orders/pkg/jwt"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/lib/pq"
)

type testDeps struct {
	cfg        *configs.Config
	db         *apidb.Db
	server     *httptest.Server
	httpClient *http.Client
}

func startTestApp(t *testing.T) *testDeps {
	t.Helper()

	// Ensure test env variables are set. If a cmd/.env exists locally, you can load it manually.
	if os.Getenv("DSN") == "" {
		os.Setenv("DSN", "postgres://postgres:my_pass@localhost:5432/orders_test?sslmode=disable")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "testsecret")
	}

	// Build app and start test server
	handler := App()
	srv := httptest.NewServer(handler)

	cfg := configs.LoadConfig()
	database := apidb.NewDb(cfg)

	return &testDeps{
		cfg:        cfg,
		db:         database,
		server:     srv,
		httpClient: srv.Client(),
	}
}

func stopTestApp(td *testDeps) {
	if td.server != nil {
		td.server.Close()
	}
}

func createTestUser(t *testing.T, db *apidb.Db, phone string) *user.User {
	t.Helper()
	repo := user.NewUserRepository(db)
	// Try find first, create if not exists
	u, err := repo.FindByPhone(phone)
	if err == nil && u != nil {
		return u
	}
	created, err := repo.Create(&user.User{Phone: phone})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	return created
}

func createTestProducts(t *testing.T, db *apidb.Db, count int) []product.Product {
	t.Helper()
	repo := product.NewProductkRepository(db)
	created := make([]product.Product, 0, count)
	for i := 0; i < count; i++ {
		p := product.NewProduct(
			"Test Product "+string(rune('A'+i)),
			"E2E description",
			pq.StringArray{"https://example.com/img.png"},
		)
		out, err := repo.Create(p)
		if err != nil {
			t.Fatalf("create product: %v", err)
		}
		created = append(created, *out)
	}
	return created
}

func deleteTestData(t *testing.T, db *apidb.Db, userID uint, productIDs []uint, orderIDs []uint) {
	t.Helper()
	// delete join rows explicitly, then orders, products, user
	if len(orderIDs) > 0 {
		if err := db.Exec("DELETE FROM order_products WHERE order_id = ANY(?)", pq.Array(orderIDs)).Error; err != nil {
			t.Fatalf("cleanup join: %v", err)
		}
		if err := db.Exec("DELETE FROM orders WHERE id = ANY(?)", pq.Array(orderIDs)).Error; err != nil {
			t.Fatalf("cleanup orders: %v", err)
		}
	}
	if len(productIDs) > 0 {
		if err := db.Exec("DELETE FROM products WHERE id = ANY(?)", pq.Array(productIDs)).Error; err != nil {
			t.Fatalf("cleanup products: %v", err)
		}
	}
	if userID != 0 {
		if err := db.Exec("DELETE FROM users WHERE id = ?", userID).Error; err != nil {
			t.Fatalf("cleanup user: %v", err)
		}
	}
}

func TestE2E_CreateOrder(t *testing.T) {
	td := startTestApp(t)
	defer stopTestApp(td)

	// Prepare data: user and products
	const phone = "79990001122"
	u := createTestUser(t, td.db, phone)
	products := createTestProducts(t, td.db, 2)
	productIDs := []uint{products[0].ID, products[1].ID}
	var createdOrderID uint
	t.Cleanup(func() {
		deleteTestData(t, td.db, u.ID, productIDs, func() []uint {
			if createdOrderID != 0 {
				return []uint{createdOrderID}
			}
			return nil
		}())
	})

	// Forge JWT for the user phone
	token, err := apijwt.NewJWT(td.cfg.Auth.Secret).Create(phone)
	if err != nil {
		t.Fatalf("make jwt: %v", err)
	}

	// Call POST /order
	body := map[string]any{"product_ids": productIDs}
	payload, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, td.server.URL+"/order", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := td.httpClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var created order.Order
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected created order id > 0")
	}
	createdOrderID = created.ID
	if created.UserID != u.ID {
		t.Fatalf("expected user_id %d, got %d", u.ID, created.UserID)
	}
	if len(created.Products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(created.Products))
	}

	// Verify GET /order/{id}
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, td.server.URL+"/order/"+strconv.FormatUint(uint64(createdOrderID), 10), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp2, err := td.httpClient.Do(req2)
	if err != nil {
		t.Fatalf("request read: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("unexpected read status: %d", resp2.StatusCode)
	}
	var fetched order.Order
	if err := json.NewDecoder(resp2.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode fetched: %v", err)
	}
	if fetched.ID != createdOrderID || fetched.UserID != u.ID || len(fetched.Products) != 2 {
		t.Fatalf("fetched order mismatch")
	}
}
