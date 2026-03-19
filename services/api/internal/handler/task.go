package handler

import (
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h *TaskHandler) GetByDate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	date := r.URL.Query().Get("date")
	tasks, err := h.task.GetByDate(r.Context(), userID, date)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to get tasks")
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{"tasks": tasks})
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Title string `json:"title"`
		Date  string `json:"date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	task, err := h.task.Create(r.Context(), userID, req.Title, req.Date)
	if err != nil {
		if err == service.ErrValidation {
			ErrJSON(w, http.StatusBadRequest, "title required")
			return
		}
		ErrJSON(w, http.StatusInternalServerError, "failed to create task")
		return
	}
	JSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid task id")
		return
	}
	cur, err := h.task.GetByID(r.Context(), userID, taskID)
	if err != nil || cur == nil {
		ErrJSON(w, http.StatusNotFound, "task not found")
		return
	}
	var req struct {
		Title     string `json:"title"`
		Date      string `json:"date"`
		Completed *bool  `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid body")
		return
	}
	title, date, completed := string(cur.Type), cur.Date, cur.Completed
	if req.Title != "" {
		title = req.Title
	}
	if req.Date != "" {
		date = req.Date
	}
	if req.Completed != nil {
		completed = *req.Completed
	}
	task, err := h.task.Update(r.Context(), userID, taskID, title, date, completed)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to update task")
		return
	}
	JSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		ErrJSON(w, http.StatusBadRequest, "invalid task id")
		return
	}
	if err := h.task.Delete(r.Context(), userID, taskID); err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to delete task")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
