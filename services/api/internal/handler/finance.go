package handler

import (
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/model"
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type FinanceHandler struct {
	finance *service.FinanceService
}

func NewFinanceHandler(finance *service.FinanceService) *FinanceHandler {
	return &FinanceHandler{finance: finance}
}

func (h *FinanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Amount float64 `json:"amount"`
		Type   string  `json:"type"`
		Note   string  `json:"note"`
		Date   string  `json:"date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Date == "" {
		ErrJSON(w, http.StatusBadRequest, "date required (YYYY-MM-DD)")
		return
	}
	financeType := parseFinanceType(req.Type)
	fin, err := h.finance.Create(r.Context(), userID, req.Amount, financeType, req.Note, req.Date)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to create finance entry")
		return
	}
	JSON(w, http.StatusCreated, fin)
}

func (h *FinanceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.finance.List(r.Context(), userID)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to list finances")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{"finances": list})
}

func parseFinanceType(s string) model.FinanceType {
	switch s {
	case "freelance":
		return model.FinanceFreelance
	case "other":
		return model.FinanceOther
	default:
		return model.FinanceSalary
	}
}
