package categories

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/montra/montra/services/api/internal/auth"
)

type mockCategoryRepo struct {
	categories []Category
}

func (m *mockCategoryRepo) List(ctx context.Context, userID string, txType string) ([]Category, error) {
	var res []Category
	for _, c := range m.categories {
		if c.UserID == nil || *c.UserID == userID {
			if txType == "" || c.TransactionType == txType {
				res = append(res, c)
			}
		}
	}
	return res, nil
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, userID, categoryID string) (*Category, error) {
	for _, c := range m.categories {
		if c.ID == categoryID && (c.UserID == nil || *c.UserID == userID) {
			return &c, nil
		}
	}
	return nil, ErrCategoryNotFound
}

func (m *mockCategoryRepo) Create(ctx context.Context, c *Category) error {
	c.ID = "cat-custom-123"
	c.IsSystem = false
	c.IsArchived = false
	c.CreatedAt = time.Now().UTC()
	c.UpdatedAt = time.Now().UTC()
	m.categories = append(m.categories, *c)
	return nil
}

func TestCategoryEndpoints(t *testing.T) {
	repo := &mockCategoryRepo{
		categories: []Category{
			{
				ID:              "cat-sys-1",
				Name:            "Makanan & Minuman",
				TransactionType: "expense",
				IsSystem:        true,
			},
			{
				ID:              "cat-sys-2",
				Name:            "Gaji",
				TransactionType: "income",
				IsSystem:        true,
			},
		},
	}

	tm, _ := auth.NewTokenManager("test-jwt-secret-with-at-least-32-characters", 1*time.Hour)
	svc := NewService(repo)
	handler := NewHandler(svc, tm)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	userToken, _ := tm.Generate("usr-uuid-1", "user@montra.id")

	// 1. List all categories with token
	reqList := httptest.NewRequest("GET", "/api/v1/categories", nil)
	reqList.Header.Set("Authorization", "Bearer "+userToken)
	wList := httptest.NewRecorder()
	mux.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list categories, got %d", wList.Code)
	}

	var listResp struct {
		Data []CategoryDTO `json:"data"`
	}
	if err := json.NewDecoder(wList.Body).Decode(&listResp); err != nil {
		t.Fatalf("failed to parse list response: %v", err)
	}
	if len(listResp.Data) != 2 {
		t.Errorf("expected 2 categories, got %d", len(listResp.Data))
	}

	// 2. Filter by type=expense
	reqFilter := httptest.NewRequest("GET", "/api/v1/categories?type=expense", nil)
	reqFilter.Header.Set("Authorization", "Bearer "+userToken)
	wFilter := httptest.NewRecorder()
	mux.ServeHTTP(wFilter, reqFilter)

	var filterResp struct {
		Data []CategoryDTO `json:"data"`
	}
	_ = json.NewDecoder(wFilter.Body).Decode(&filterResp)
	if len(filterResp.Data) != 1 || filterResp.Data[0].Name != "Makanan & Minuman" {
		t.Errorf("expected only Makanan & Minuman for expense, got %v", filterResp.Data)
	}

	// 3. Create custom category
	createBody := map[string]string{
		"name":             "Kopi & Cemilan",
		"transaction_type": "expense",
		"icon":             "coffee",
		"color":            "#795548",
	}
	jsonBody, _ := json.Marshal(createBody)
	reqCreate := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewReader(jsonBody))
	reqCreate.Header.Set("Authorization", "Bearer "+userToken)
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	mux.ServeHTTP(wCreate, reqCreate)

	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on category creation, got %d: %s", wCreate.Code, wCreate.Body.String())
	}

	// 4. Invalid transaction_type returns 422
	badBody := map[string]string{
		"name":             "Invalid Type",
		"transaction_type": "investment",
	}
	badJSON, _ := json.Marshal(badBody)
	reqBad := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewReader(badJSON))
	reqBad.Header.Set("Authorization", "Bearer "+userToken)
	reqBad.Header.Set("Content-Type", "application/json")
	wBad := httptest.NewRecorder()
	mux.ServeHTTP(wBad, reqBad)

	if wBad.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 Unprocessable Entity, got %d", wBad.Code)
	}
}
