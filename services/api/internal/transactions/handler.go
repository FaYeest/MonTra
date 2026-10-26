package transactions

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/montra/montra/services/api/internal/auth"
)

// Handler handles HTTP requests for transactions.
type Handler struct {
	service *Service
	tm      *auth.TokenManager
}

// NewHandler constructs a Transaction HTTP Handler.
func NewHandler(service *Service, tm *auth.TokenManager) *Handler {
	return &Handler{
		service: service,
		tm:      tm,
	}
}

// RegisterRoutes registers transaction endpoints on the provided router.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	authMiddleware := auth.RequireAuth(h.tm)

	mux.Handle("GET /api/v1/transactions", authMiddleware(http.HandlerFunc(h.HandleList)))
	mux.Handle("POST /api/v1/transactions", authMiddleware(http.HandlerFunc(h.HandleCreate)))
	mux.Handle("GET /api/v1/transactions/{id}", authMiddleware(http.HandlerFunc(h.HandleGetByID)))
	mux.Handle("PATCH /api/v1/transactions/{id}", authMiddleware(http.HandlerFunc(h.HandleUpdate)))
	mux.Handle("DELETE /api/v1/transactions/{id}", authMiddleware(http.HandlerFunc(h.HandleDelete)))
}

// HandleList handles GET /api/v1/transactions
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	filter := Filter{
		AccountID:  q.Get("account_id"),
		CategoryID: q.Get("category_id"),
		Type:       q.Get("type"),
		Search:     q.Get("search"),
		Limit:      limit,
		Offset:     offset,
	}

	if q.Get("start_date") != "" {
		if t, err := parseDate(q.Get("start_date")); err == nil {
			filter.StartDate = &t
		}
	}
	if q.Get("end_date") != "" {
		if t, err := parseDate(q.Get("end_date")); err == nil {
			filter.EndDate = &t
		}
	}

	items, total, err := h.service.ListTransactions(r.Context(), userID, filter)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat transaksi.")
		return
	}

	if items == nil {
		items = []TransactionDetailDTO{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": map[string]any{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// HandleCreate handles POST /api/v1/transactions
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	// Honor Idempotency-Key header if not present in payload
	if (req.IdempotencyKey == nil || *req.IdempotencyKey == "") && r.Header.Get("Idempotency-Key") != "" {
		headerKey := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		req.IdempotencyKey = &headerKey
	}

	item, err := h.service.CreateTransaction(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAmount), errors.Is(err, ErrInvalidType),
			errors.Is(err, ErrInvalidDescription), errors.Is(err, ErrCategoryNotValid),
			errors.Is(err, ErrInvalidDate):
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, ErrAccountNotFound):
			writeJSONError(w, http.StatusNotFound, "ACCOUNT_NOT_FOUND", err.Error())
		default:
			slog.Error("failed to create transaction", "user_id", userID, "error", err)
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mencatat transaksi.")
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"data": item,
	})
}

// HandleGetByID handles GET /api/v1/transactions/{id}
func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	txID := r.PathValue("id")
	if txID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID transaksi tidak boleh kosong.")
		return
	}

	item, err := h.service.GetTransaction(r.Context(), userID, txID)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Transaksi tidak ditemukan.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil data transaksi.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": item,
	})
}

// HandleUpdate handles PATCH /api/v1/transactions/{id}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	txID := r.PathValue("id")
	if txID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID transaksi tidak boleh kosong.")
		return
	}

	var req UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	item, err := h.service.UpdateTransaction(r.Context(), userID, txID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTransactionNotFound):
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Transaksi tidak ditemukan.")
		case errors.Is(err, ErrInvalidAmount), errors.Is(err, ErrInvalidType),
			errors.Is(err, ErrInvalidDescription), errors.Is(err, ErrCategoryNotValid),
			errors.Is(err, ErrInvalidDate):
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, ErrAccountNotFound):
			writeJSONError(w, http.StatusNotFound, "ACCOUNT_NOT_FOUND", err.Error())
		default:
			slog.Error("failed to update transaction", "user_id", userID, "error", err)
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memperbarui transaksi.")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": item,
	})
}

// HandleDelete handles DELETE /api/v1/transactions/{id}
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	txID := r.PathValue("id")
	if txID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID transaksi tidak boleh kosong.")
		return
	}

	if err := h.service.DeleteTransaction(r.Context(), userID, txID); err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Transaksi tidak ditemukan.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menghapus transaksi.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]string{
			"message": "Transaksi berhasil dihapus.",
		},
	})
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
