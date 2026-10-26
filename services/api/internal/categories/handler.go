package categories

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/montra/montra/services/api/internal/auth"
)

// Handler handles HTTP requests for categories.
type Handler struct {
	service *Service
	tm      *auth.TokenManager
}

// NewHandler constructs a Category HTTP Handler.
func NewHandler(service *Service, tm *auth.TokenManager) *Handler {
	return &Handler{
		service: service,
		tm:      tm,
	}
}

// RegisterRoutes attaches category routes to the router with authentication middleware.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	authMiddleware := auth.RequireAuth(h.tm)

	mux.Handle("GET /api/v1/categories", authMiddleware(http.HandlerFunc(h.HandleList)))
	mux.Handle("POST /api/v1/categories", authMiddleware(http.HandlerFunc(h.HandleCreate)))
}

// HandleList handles GET /api/v1/categories?type=income|expense
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	txType := r.URL.Query().Get("type")
	cats, err := h.service.ListCategories(r.Context(), userID, txType)
	if err != nil {
		if errors.Is(err, ErrInvalidTransactionType) {
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat kategori.")
		return
	}

	writeJSONData(w, http.StatusOK, cats)
}

// HandleCreate handles POST /api/v1/categories
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	cat, err := h.service.CreateCategory(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCategoryName), errors.Is(err, ErrInvalidTransactionType):
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuat kategori.")
		}
		return
	}

	writeJSONData(w, http.StatusCreated, cat)
}

func writeJSONData(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": data,
	})
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
