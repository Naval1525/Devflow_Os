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

type OpportunityHandler struct {
	opp *service.OpportunityService
}

func NewOpportunityHandler(opp *service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{opp: opp}
}

func (h *OpportunityHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Stage  string `json:"stage"`
		Source string `json:"source"`
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Name == "" {
		ErrJSON(w, http.StatusBadRequest, "name required")
		return
	}
	oppType := parseOppType(req.Type)
	stage := parseOppStage(req.Stage)
	opp, err := h.opp.Create(r.Context(), userID, req.Name, oppType, stage, req.Source, req.Notes)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to create opportunity")
		return
	}
	JSON(w, http.StatusCreated, opp)
}

func (h *OpportunityHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.opp.List(r.Context(), userID)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to list opportunities")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{"opportunities": list})
}

func (h *OpportunityHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		Name   *string `json:"name"`
		Type   *string `json:"type"`
		Stage  *string `json:"stage"`
		Source *string `json:"source"`
		Notes  *string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	var oppType *model.OpportunityType
	var stage *model.OpportunityStage
	if req.Type != nil {
		t := parseOppType(*req.Type)
		oppType = &t
	}
	if req.Stage != nil {
		s := parseOppStage(*req.Stage)
		stage = &s
	}
	opp, err := h.opp.Update(r.Context(), id, userID, req.Name, req.Source, req.Notes, oppType, stage)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to update opportunity")
		return
	}
	JSON(w, http.StatusOK, opp)
}

func parseOppType(s string) model.OpportunityType {
	if s == "freelance" {
		return model.OppTypeFreelance
	}
	return model.OppTypeJob
}

func parseOppStage(s string) model.OpportunityStage {
	switch s {
	case "interview":
		return model.OppStageInterview
	case "closed":
		return model.OppStageClosed
	default:
		return model.OppStageApplied
	}
}
