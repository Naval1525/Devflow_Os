package handler

import (
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Signup is disabled; registration is not allowed.
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ErrJSON(w, http.StatusForbidden, "Access denied")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Email == "" || req.Password == "" {
		ErrJSON(w, http.StatusBadRequest, "email and password required")
		return
	}
	token, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrAccessDenied {
			ErrJSON(w, http.StatusForbidden, "Access denied")
			return
		}
		if err == service.ErrInvalidCredentials {
			ErrJSON(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		ErrJSON(w, http.StatusInternalServerError, "login failed")
		return
	}
	JSON(w, http.StatusOK, map[string]string{"token": token})
}
