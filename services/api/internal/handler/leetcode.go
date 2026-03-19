package handler

import (
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/model"
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type LeetCodeHandler struct {
	leetcode *service.LeetCodeService
}

func NewLeetCodeHandler(leetcode *service.LeetCodeService) *LeetCodeHandler {
	return &LeetCodeHandler{leetcode: leetcode}
}

func (h *LeetCodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		ProblemName string `json:"problem_name"`
		Difficulty  string `json:"difficulty"`
		Approach    string `json:"approach"`
		Mistake     string `json:"mistake"`
		TimeTaken   *int   `json:"time_taken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.ProblemName == "" {
		ErrJSON(w, http.StatusBadRequest, "problem_name required")
		return
	}
	diff := parseDifficulty(req.Difficulty)
	log, err := h.leetcode.Create(r.Context(), userID, req.ProblemName, diff, req.Approach, req.Mistake, req.TimeTaken)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to create log")
		return
	}
	JSON(w, http.StatusCreated, log)
}

func (h *LeetCodeHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.leetcode.List(r.Context(), userID)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to list logs")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{"logs": list})
}

func parseDifficulty(s string) model.Difficulty {
	switch s {
	case "medium":
		return model.DifficultyMedium
	case "hard":
		return model.DifficultyHard
	default:
		return model.DifficultyEasy
	}
}
