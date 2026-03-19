package service

import (
	"context"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type IdeaService struct {
	ideaRepo IdeaRepository
}

type IdeaRepository interface {
	Create(ctx context.Context, userID uuid.UUID, hook, idea string, ideaType model.IdeaType, status model.IdeaStatus) (*model.Idea, error)
	List(ctx context.Context, userID uuid.UUID, filterType *model.IdeaType, filterStatus *model.IdeaStatus) ([]*model.Idea, error)
	Update(ctx context.Context, id, userID uuid.UUID, hook, idea *string, ideaType *model.IdeaType, status *model.IdeaStatus) (*model.Idea, error)
}

func NewIdeaService(ideaRepo IdeaRepository) *IdeaService {
	return &IdeaService{ideaRepo: ideaRepo}
}

func (s *IdeaService) Create(ctx context.Context, userID uuid.UUID, hook, idea string, ideaType model.IdeaType, status model.IdeaStatus) (*model.Idea, error) {
	if status == "" {
		status = model.IdeaStatusIdea
	}
	return s.ideaRepo.Create(ctx, userID, hook, idea, ideaType, status)
}

func (s *IdeaService) List(ctx context.Context, userID uuid.UUID, filterType *model.IdeaType, filterStatus *model.IdeaStatus) ([]*model.Idea, error) {
	return s.ideaRepo.List(ctx, userID, filterType, filterStatus)
}

func (s *IdeaService) Update(ctx context.Context, id, userID uuid.UUID, hook, idea *string, ideaType *model.IdeaType, status *model.IdeaStatus) (*model.Idea, error) {
	return s.ideaRepo.Update(ctx, id, userID, hook, idea, ideaType, status)
}
