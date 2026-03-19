package service

import (
	"context"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type OpportunityService struct {
	oppRepo OpportunityRepository
}

type OpportunityRepository interface {
	Create(ctx context.Context, userID uuid.UUID, name string, oppType model.OpportunityType, stage model.OpportunityStage, source, notes string) (*model.Opportunity, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.Opportunity, error)
	Update(ctx context.Context, id, userID uuid.UUID, name, source, notes *string, oppType *model.OpportunityType, stage *model.OpportunityStage) (*model.Opportunity, error)
}

func NewOpportunityService(oppRepo OpportunityRepository) *OpportunityService {
	return &OpportunityService{oppRepo: oppRepo}
}

func (s *OpportunityService) Create(ctx context.Context, userID uuid.UUID, name string, oppType model.OpportunityType, stage model.OpportunityStage, source, notes string) (*model.Opportunity, error) {
	return s.oppRepo.Create(ctx, userID, name, oppType, stage, source, notes)
}

func (s *OpportunityService) List(ctx context.Context, userID uuid.UUID) ([]*model.Opportunity, error) {
	return s.oppRepo.List(ctx, userID)
}

func (s *OpportunityService) Update(ctx context.Context, id, userID uuid.UUID, name, source, notes *string, oppType *model.OpportunityType, stage *model.OpportunityStage) (*model.Opportunity, error) {
	return s.oppRepo.Update(ctx, id, userID, name, source, notes, oppType, stage)
}
