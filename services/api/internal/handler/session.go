package handler

import (
	"devflowos/api/internal/middleware"
	"devflowos/api/internal/service"
	"net/http"

	"github.com/google/uuid"
)

type SessionHandler struct {
	session *service.SessionService
}

func NewSessionHandler(session *service.SessionService) *SessionHandler {
	return &SessionHandler{session: session}
}

func (h *SessionHandler) Start(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	sess, err := h.session.Start(r.Context(), userID)
	if err != nil {
		if err == service.ErrActiveSessionExists {
			ErrJSON(w, http.StatusConflict, "already have an active session")
			return
		}
		ErrJSON(w, http.StatusInternalServerError, "failed to start session")
		return
	}
	JSON(w, http.StatusOK, sess)
}

func (h *SessionHandler) End(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	sess, err := h.session.End(r.Context(), userID)
	if err != nil {
		ErrJSON(w, http.StatusBadRequest, "no active session")
		return
	}
	JSON(w, http.StatusOK, sess)
}

func (h *SessionHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == uuid.Nil {
		ErrJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	sess, err := h.session.GetActive(r.Context(), userID)
	if err != nil {
		ErrJSON(w, http.StatusInternalServerError, "failed to get session")
		return
	}
	if sess == nil {
		JSON(w, http.StatusOK, nil)
		return
	}
	JSON(w, http.StatusOK, sess)
}
