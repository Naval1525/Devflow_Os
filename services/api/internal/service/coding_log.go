package service

import (
	"context"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type CodingLogService struct {
	repo CodingLogRepository
}

type CodingLogRepository interface {
	Create(ctx context.Context, userID uuid.UUID, title, description string) (*model.CodingLog, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.CodingLog, error)
}

func NewCodingLogService(repo CodingLogRepository) *CodingLogService {
	return &CodingLogService{repo: repo}
}

func (s *CodingLogService) Create(ctx context.Context, userID uuid.UUID, title, description string) (*model.CodingLog, error) {
	return s.repo.Create(ctx, userID, title, description)
}

func (s *CodingLogService) List(ctx context.Context, userID uuid.UUID) ([]*model.CodingLog, error) {
	return s.repo.List(ctx, userID)
}
