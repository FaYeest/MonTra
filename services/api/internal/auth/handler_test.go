package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/montra/montra/services/api/internal/users"
)

type mockUserRepo struct {
	usersByEmail map[string]*users.User
	usersByID    map[string]*users.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		usersByEmail: make(map[string]*users.User),
		usersByID:    make(map[string]*users.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, u *users.User) error {
	if _, exists := m.usersByEmail[u.Email]; exists {
		return users.ErrEmailAlreadyExists
	}
	u.ID = "usr-generated-123"
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()
	m.usersByEmail[u.Email] = u
	m.usersByID[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	u, ok := m.usersByEmail[email]
	if !ok {
		return nil, users.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*users.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, users.ErrNotFound
	}
	return u, nil
}

func setupTestRouter() (http.Handler, *mockUserRepo, *TokenManager) {
	repo := newMockUserRepo()
	hasher := NewHasher()
	tm, _ := NewTokenManager("test-jwt-secret-with-at-least-32-characters", 2*time.Hour)
	svc := NewService(repo, hasher, tm)
	h := NewHandler(svc, tm)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux, repo, tm
}

func TestAuthEndpoints(t *testing.T) {
	router, _, tm := setupTestRouter()

	// 1. Register user
	regPayload := map[string]string{
		"email":        "farras@example.com",
		"password":     "superSecret123",
		"display_name": "Farras",
	}
	body, _ := json.Marshal(regPayload)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on register, got %d: %s", w.Code, w.Body.String())
	}

	var regResp struct {
		Data AuthResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&regResp); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}
	token := regResp.Data.Token
	if token == "" {
		t.Fatalf("expected token in register response, got empty")
	}
	if regResp.Data.User.Email != "farras@example.com" {
		t.Errorf("expected email farras@example.com, got %s", regResp.Data.User.Email)
	}

	// 2. Duplicate registration should return 409 Conflict
	wDup := httptest.NewRecorder()
	reqDup := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	reqDup.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wDup, reqDup)
	if wDup.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict on duplicate register, got %d", wDup.Code)
	}

	// 3. Login with correct credentials
	loginPayload := map[string]string{
		"email":    "farras@example.com",
		"password": "superSecret123",
	}
	loginBody, _ := json.Marshal(loginPayload)
	wLogin := httptest.NewRecorder()
	reqLogin := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wLogin, reqLogin)
	if wLogin.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on login, got %d: %s", wLogin.Code, wLogin.Body.String())
	}

	// 4. Login with incorrect password
	badLoginPayload := map[string]string{
		"email":    "farras@example.com",
		"password": "wrongPassword123",
	}
	badBody, _ := json.Marshal(badLoginPayload)
	wBadLogin := httptest.NewRecorder()
	reqBadLogin := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(badBody))
	reqBadLogin.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wBadLogin, reqBadLogin)
	if wBadLogin.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized on bad credentials, got %d", wBadLogin.Code)
	}

	// 5. Access /api/v1/me without token -> 401
	wNoAuth := httptest.NewRecorder()
	reqNoAuth := httptest.NewRequest("GET", "/api/v1/me", nil)
	router.ServeHTTP(wNoAuth, reqNoAuth)
	if wNoAuth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized without token, got %d", wNoAuth.Code)
	}

	// 6. Access /api/v1/me with valid Bearer token -> 200
	wMe := httptest.NewRecorder()
	reqMe := httptest.NewRequest("GET", "/api/v1/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(wMe, reqMe)
	if wMe.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on /me, got %d: %s", wMe.Code, wMe.Body.String())
	}

	// 7. Access /api/v1/me with tampered token -> 401
	wTampered := httptest.NewRecorder()
	reqTampered := httptest.NewRequest("GET", "/api/v1/me", nil)
	reqTampered.Header.Set("Authorization", "Bearer invalid.jwt.token")
	router.ServeHTTP(wTampered, reqTampered)
	if wTampered.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized on tampered token, got %d", wTampered.Code)
	}
	_ = tm
}
