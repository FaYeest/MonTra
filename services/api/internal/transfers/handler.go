package transfers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/montra/montra/services/api/internal/auth"
)

// Handler handles HTTP requests for fund transfers.
type Handler struct {
	service *Service
	tm      *auth.TokenManager
}

// NewHandler constructs a Transfer HTTP Handler.
func NewHandler(service *Service, tm *auth.TokenManager) *Handler {
	return &Handler{
		service: service,
		tm:      tm,
	}
}

// RegisterRoutes attaches transfer routes to the router.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	authMiddleware := auth.RequireAuth(h.tm)

	mux.Handle("GET /api/v1/transfers", authMiddleware(http.HandlerFunc(h.HandleList)))
	mux.Handle("POST /api/v1/transfers", authMiddleware(http.HandlerFunc(h.HandleCreate)))
	mux.Handle("GET /api/v1/transfers/{id}", authMiddleware(http.HandlerFunc(h.HandleGetByID)))
	mux.Handle("DELETE /api/v1/transfers/{id}", authMiddleware(http.HandlerFunc(h.HandleDelete)))
}

// HandleList handles GET /api/v1/transfers
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	items, total, err := h.service.ListTransfers(r.Context(), userID, limit, offset)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat daftar transfer.")
		return
	}

	if items == nil {
		items = []TransferDetailDTO{}
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

// HandleCreate handles POST /api/v1/transfers
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	var req CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	if (req.IdempotencyKey == nil || *req.IdempotencyKey == "") && r.Header.Get("Idempotency-Key") != "" {
		headerKey := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		req.IdempotencyKey = &headerKey
	}

	item, err := h.service.CreateTransfer(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrSameAccountTransfer), errors.Is(err, ErrInvalidAmount), errors.Is(err, ErrInvalidFee):
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, ErrAccountNotFound):
			writeJSONError(w, http.StatusNotFound, "ACCOUNT_NOT_FOUND", err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mencatat transfer.")
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"data": item,
	})
}

// HandleGetByID handles GET /api/v1/transfers/{id}
func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	trID := r.PathValue("id")
	if trID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID transfer tidak boleh kosong.")
		return
	}

	item, err := h.service.GetTransfer(r.Context(), userID, trID)
	if err != nil {
		if errors.Is(err, ErrTransferNotFound) {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Transfer tidak ditemukan.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil data transfer.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": item,
	})
}

// HandleDelete handles DELETE /api/v1/transfers/{id}
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	trID := r.PathValue("id")
	if trID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID transfer tidak boleh kosong.")
		return
	}

	if err := h.service.DeleteTransfer(r.Context(), userID, trID); err != nil {
		if errors.Is(err, ErrTransferNotFound) {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Transfer tidak ditemukan.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menghapus data transfer.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]string{
			"message": "Transfer berhasil dihapus.",
		},
	})
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
