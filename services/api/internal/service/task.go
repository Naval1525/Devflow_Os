package service

import (
	"context"
	"devflowos/api/internal/model"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrValidation = errors.New("validation error")

type TaskService struct {
	taskRepo TaskRepository
}

type TaskRepository interface {
	GetByDate(ctx context.Context, userID uuid.UUID, date string) ([]*model.Task, error)
	Create(ctx context.Context, userID uuid.UUID, title string, date string) (*model.Task, error)
	UpdateByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, title string, date string, completed bool) error
	DeleteByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error
}

func NewTaskService(taskRepo TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) GetToday(ctx context.Context, userID uuid.UUID) ([]*model.Task, error) {
	date := time.Now().UTC().Format("2006-01-02")
	return s.taskRepo.GetByDate(ctx, userID, date)
}

func (s *TaskService) GetByDate(ctx context.Context, userID uuid.UUID, date string) ([]*model.Task, error) {
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	return s.taskRepo.GetByDate(ctx, userID, date)
}

func (s *TaskService) Create(ctx context.Context, userID uuid.UUID, title string, date string) (*model.Task, error) {
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	if title == "" {
		return nil, ErrValidation
	}
	return s.taskRepo.Create(ctx, userID, title, date)
}

func (s *TaskService) Update(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, title string, date string, completed bool) (*model.Task, error) {
	return s.taskRepo.UpdateByID(ctx, userID, taskID, title, date, completed)
}

func (s *TaskService) Delete(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {
	return s.taskRepo.DeleteByID(ctx, userID, taskID)
}
