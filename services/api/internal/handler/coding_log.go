package handler

import (
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type CodingLogHandler struct {
	svc *service.CodingLogService
}

func NewCodingLogHandler(svc *service.CodingLogService) *CodingLogHandler {
	return &CodingLogHandler{svc: svc}
}

func (h *CodingLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Title == "" {
		ErrJSON(w, http.StatusBadRequest, "title required")
		return
	}
	log, err := h.svc.Create(r.Context(), userID, req.Title, req.Description)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to create log")
		return
	}
	JSON(w, http.StatusCreated, log)
}

func (h *CodingLogHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	logs, err := h.svc.List(r.Context(), userID)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to list logs")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{"logs": logs})
}
