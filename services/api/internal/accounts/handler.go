package accounts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/montra/montra/services/api/internal/auth"
)

// Handler handles HTTP requests for financial accounts.
type Handler struct {
	service *Service
	tm      *auth.TokenManager
}

// NewHandler constructs an Account HTTP Handler.
func NewHandler(service *Service, tm *auth.TokenManager) *Handler {
	return &Handler{
		service: service,
		tm:      tm,
	}
}

// RegisterRoutes attaches account routes to the HTTP mux with auth middleware.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	authMiddleware := auth.RequireAuth(h.tm)

	mux.Handle("GET /api/v1/accounts", authMiddleware(http.HandlerFunc(h.HandleList)))
	mux.Handle("POST /api/v1/accounts", authMiddleware(http.HandlerFunc(h.HandleCreate)))
	mux.Handle("GET /api/v1/accounts/{id}", authMiddleware(http.HandlerFunc(h.HandleGetByID)))
	mux.Handle("PATCH /api/v1/accounts/{id}", authMiddleware(http.HandlerFunc(h.HandleUpdate)))
	mux.Handle("DELETE /api/v1/accounts/{id}", authMiddleware(http.HandlerFunc(h.HandleArchive)))
}

// HandleList handles GET /api/v1/accounts
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	accs, err := h.service.ListAccounts(r.Context(), userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat akun.")
		return
	}

	writeJSONData(w, http.StatusOK, accs)
}

// HandleCreate handles POST /api/v1/accounts
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	acc, err := h.service.CreateAccount(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAccountName), errors.Is(err, ErrInvalidAccountType), errors.Is(err, ErrNegativeOpeningBalance):
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuat akun.")
		}
		return
	}

	writeJSONData(w, http.StatusCreated, acc)
}

// HandleGetByID handles GET /api/v1/accounts/{id}
func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	accountID := r.PathValue("id")
	if accountID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID akun tidak boleh kosong.")
		return
	}

	acc, err := h.service.GetAccount(r.Context(), userID, accountID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Akun tidak ditemukan.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil data akun.")
		return
	}

	writeJSONData(w, http.StatusOK, acc)
}

// HandleUpdate handles PATCH /api/v1/accounts/{id}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	accountID := r.PathValue("id")
	if accountID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID akun tidak boleh kosong.")
		return
	}

	var req UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	acc, err := h.service.UpdateAccount(r.Context(), userID, accountID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrAccountNotFound):
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Akun tidak ditemukan.")
		case errors.Is(err, ErrInvalidAccountName), errors.Is(err, ErrInvalidAccountType):
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memperbarui akun.")
		}
		return
	}

	writeJSONData(w, http.StatusOK, acc)
}

// HandleArchive handles DELETE /api/v1/accounts/{id} (soft archive)
func (h *Handler) HandleArchive(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	accountID := r.PathValue("id")
	if accountID == "" {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "ID akun tidak boleh kosong.")
		return
	}

	if err := h.service.ArchiveAccount(r.Context(), userID, accountID); err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Akun tidak ditemukan.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengarsipkan akun.")
		return
	}

	writeJSONData(w, http.StatusOK, map[string]string{
		"message": "Akun berhasil diarsipkan.",
	})
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
