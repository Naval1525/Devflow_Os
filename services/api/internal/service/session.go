package service

import (
	"context"
	"devflowos/api/internal/model"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrActiveSessionExists = errors.New("already have an active session")

type SessionService struct {
	sessionRepo SessionRepository
}

type SessionRepository interface {
	Create(ctx context.Context, userID uuid.UUID, startTime time.Time) (*model.Session, error)
	GetActive(ctx context.Context, userID uuid.UUID) (*model.Session, error)
	End(ctx context.Context, id, userID uuid.UUID, endTime time.Time) (*model.Session, error)
}

func NewSessionService(sessionRepo SessionRepository) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) Start(ctx context.Context, userID uuid.UUID) (*model.Session, error) {
	active, err := s.sessionRepo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, ErrActiveSessionExists
	}
	return s.sessionRepo.Create(ctx, userID, time.Now().UTC())
}

func (s *SessionService) End(ctx context.Context, userID uuid.UUID) (*model.Session, error) {
	active, err := s.sessionRepo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return nil, errors.New("no active session")
	}
	return s.sessionRepo.End(ctx, active.ID, userID, time.Now().UTC())
}

func (s *SessionService) GetActive(ctx context.Context, userID uuid.UUID) (*model.Session, error) {
	return s.sessionRepo.GetActive(ctx, userID)
}
