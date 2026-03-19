package service

import (
	"context"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type LeetCodeService struct {
	leetcodeRepo LeetCodeRepository
}

type LeetCodeRepository interface {
	Create(ctx context.Context, userID uuid.UUID, problemName string, difficulty model.Difficulty, approach, mistake string, timeTaken *int) (*model.LeetCodeLog, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.LeetCodeLog, error)
}

func NewLeetCodeService(leetcodeRepo LeetCodeRepository) *LeetCodeService {
	return &LeetCodeService{leetcodeRepo: leetcodeRepo}
}

func (s *LeetCodeService) Create(ctx context.Context, userID uuid.UUID, problemName string, difficulty model.Difficulty, approach, mistake string, timeTaken *int) (*model.LeetCodeLog, error) {
	return s.leetcodeRepo.Create(ctx, userID, problemName, difficulty, approach, mistake, timeTaken)
}

func (s *LeetCodeService) List(ctx context.Context, userID uuid.UUID) ([]*model.LeetCodeLog, error) {
	return s.leetcodeRepo.List(ctx, userID)
}
