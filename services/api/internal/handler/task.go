package handler

import (
	"devflowos/api/internal/model"
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"

	"devflowos/api/internal/middleware"

	"github.com/google/uuid"
)

type TaskHandler struct {
	task *service.TaskService
}

func NewTaskHandler(task *service.TaskService) *TaskHandler {
	return &TaskHandler{task: task}
}

func (h *TaskHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	tasks, err := h.task.GetToday(r.Context(), userID)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to get tasks")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{"tasks": tasks})
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Type      string `json:"type"`
		Date      string `json:"date"`
		Completed bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	var taskType model.TaskType
	switch req.Type {
	case "coding":
		taskType = model.TaskCoding
	case "leetcode":
		taskType = model.TaskLeetcode
	case "content":
		taskType = model.TaskContent
	default:
		ErrJSON(w, http.StatusBadRequest, "type must be coding, leetcode, or content")
		return
	}
	task, err := h.task.Update(r.Context(), userID, taskType, req.Date, req.Completed)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to update task")
		return
	}
	JSON(w, http.StatusOK, task)
}

