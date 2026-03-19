package handler

import (
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/model"
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type IdeaHandler struct {
	idea *service.IdeaService
}

func NewIdeaHandler(idea *service.IdeaService) *IdeaHandler {
	return &IdeaHandler{idea: idea}
}

func (h *IdeaHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Hook   string `json:"hook"`
		Idea   string `json:"idea"`
		Type   string `json:"type"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Hook == "" {
		ErrJSON(w, http.StatusBadRequest, "hook required")
		return
	}
	ideaType := parseIdeaType(req.Type)
	status := parseIdeaStatus(req.Status)
	idea, err := h.idea.Create(r.Context(), userID, req.Hook, req.Idea, ideaType, status)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to create idea")
		return
	}
	JSON(w, http.StatusCreated, idea)
}

func (h *IdeaHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var typeFilter *model.IdeaType
	if t := r.URL.Query().Get("type"); t != "" {
		tt := parseIdeaType(t)
		typeFilter = &tt
	}
	var statusFilter *model.IdeaStatus
	if s := r.URL.Query().Get("status"); s != "" {
		ss := parseIdeaStatus(s)
		statusFilter = &ss
	}
	list, err := h.idea.List(r.Context(), userID, typeFilter, statusFilter)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to list ideas")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{"ideas": list})
}

func (h *IdeaHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Hook   *string `json:"hook"`
		Idea   *string `json:"idea"`
		Type   *string `json:"type"`
		Status *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	var ideaType *model.IdeaType
	var status *model.IdeaStatus
	if req.Type != nil {
		t := parseIdeaType(*req.Type)
		ideaType = &t
	}
	if req.Status != nil {
		s := parseIdeaStatus(*req.Status)
		status = &s
	}
	idea, err := h.idea.Update(r.Context(), id, userID, req.Hook, req.Idea, ideaType, status)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to update idea")
		return
	}
	JSON(w, http.StatusOK, idea)
}

func parseIdeaType(s string) model.IdeaType {
	switch s {
	case "reel":
		return model.IdeaTypeReel
	case "thread":
		return model.IdeaTypeThread
	case "linkedin":
		return model.IdeaTypeLinkedin
	default:
		return model.IdeaTypeTweet
	}
}

func parseIdeaStatus(s string) model.IdeaStatus {
	switch s {
	case "ready":
		return model.IdeaStatusReady
	case "posted":
		return model.IdeaStatusPosted
	default:
		return model.IdeaStatusIdea
	}
}
