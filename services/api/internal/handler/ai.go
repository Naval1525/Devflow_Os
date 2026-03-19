package handler

import (
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"
)

type AIHandler struct {
	ai *service.AIContentService
}

func NewAIHandler(ai *service.AIContentService) *AIHandler {
	return &AIHandler{ai: ai}
}

func (h *AIHandler) GenerateContent(w http.ResponseWriter, r *http.Request) {
	_ = middleware.UserIDFromContext(r.Context())
	var req service.GenerateContentInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Text == "" {
		ErrJSON(w, http.StatusBadRequest, "text required")
		return
	}
	out, err := h.ai.Generate(r.Context(), req)
	if err != nil {
		if err == service.ErrGeminiNotConfigured {
			ErrJSON(w, http.StatusServiceUnavailable, "AI not configured")
			return
		}
		ErrJSON(w, http.StatusInternalServerError, "generation failed")
		return
	}
	JSON(w, http.StatusOK, out)
}
