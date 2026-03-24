package service

import (
	"context"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type FinanceService struct {
	financeRepo FinanceRepository
}

type FinanceRepository interface {
	Create(ctx context.Context, userID uuid.UUID, amount float64, financeType model.FinanceType, note, date string) (*model.Finance, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.Finance, error)
	Delete(ctx context.Context, userID, financeID uuid.UUID) error
	SumByMonth(ctx context.Context, userID uuid.UUID, year, month int) (float64, error)
}

func NewFinanceService(financeRepo FinanceRepository) *FinanceService {
	return &FinanceService{financeRepo: financeRepo}
}

func (s *FinanceService) Create(ctx context.Context, userID uuid.UUID, amount float64, financeType model.FinanceType, note, date string) (*model.Finance, error) {
	return s.financeRepo.Create(ctx, userID, amount, financeType, note, date)
}

func (s *FinanceService) List(ctx context.Context, userID uuid.UUID) ([]*model.Finance, error) {
	return s.financeRepo.List(ctx, userID)
}

func (s *FinanceService) Delete(ctx context.Context, userID, financeID uuid.UUID) error {
	return s.financeRepo.Delete(ctx, userID, financeID)
}

func (s *FinanceService) SumByMonth(ctx context.Context, userID uuid.UUID, year, month int) (float64, error) {
	return s.financeRepo.SumByMonth(ctx, userID, year, month)
}
