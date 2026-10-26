package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/montra/montra/services/api/internal/users"
)

// Handler exposes HTTP handlers for authentication and user profile.
type Handler struct {
	service *Service
	tm      *TokenManager
}

// NewHandler constructs an auth HTTP Handler.
func NewHandler(service *Service, tm *TokenManager) *Handler {
	return &Handler{
		service: service,
		tm:      tm,
	}
}

// RegisterRoutes attaches authentication endpoints to the provided HTTP router.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", h.HandleRegister)
	mux.HandleFunc("POST /api/v1/auth/login", h.HandleLogin)

	// Protected routes wrapped with RequireAuth
	authMiddleware := RequireAuth(h.tm)
	mux.Handle("GET /api/v1/me", authMiddleware(http.HandlerFunc(h.HandleMe)))
}

// HandleRegister handles user registration.
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrEmailAlreadyExists):
			writeJSONError(w, http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email sudah terdaftar.")
		case errors.Is(err, ErrInvalidEmail), errors.Is(err, ErrInvalidDisplayName),
			errors.Is(err, ErrPasswordTooShort), errors.Is(err, ErrPasswordTooLong):
			writeJSONError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mendaftarkan akun.")
		}
		return
	}

	writeJSONData(w, http.StatusCreated, resp)
}

// HandleLogin handles user authentication.
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", "Format payload JSON tidak valid.")
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeJSONError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email atau kata sandi tidak valid.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memproses autentikasi.")
		return
	}

	writeJSONData(w, http.StatusOK, resp)
}

// HandleMe returns the profile of the authenticated user.
func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Autentikasi diperlukan.")
		return
	}

	profile, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "USER_NOT_FOUND", "Pengguna tidak ditemukan.")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil data profil.")
		return
	}

	writeJSONData(w, http.StatusOK, map[string]any{
		"user": profile,
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
