package handler

import (
	"devflowos/api/internal/model"
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

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
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
	u, token, err := h.auth.Signup(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrEmailExists {
			ErrJSON(w, http.StatusConflict, "email already registered")
			return
		}
		ErrJSON(w, http.StatusInternalServerError, "signup failed")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{
		"user":  userResponse(u),
		"token": token,
	})
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
		if err == service.ErrInvalidCredentials {
			ErrJSON(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		ErrJSON(w, http.StatusInternalServerError, "login failed")
		return
	}
	JSON(w, http.StatusOK, map[string]string{"token": token})
}

func userResponse(u *model.User) map[string]interface{} {
	return map[string]interface{}{
		"id":    u.ID.String(),
		"email": u.Email,
	}
}
