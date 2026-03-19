package service

import (
	"context"
	"devflowos/api/internal/model"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo TaskRepository
}

type TaskRepository interface {
	GetToday(ctx context.Context, userID uuid.UUID, date string) ([]*model.Task, error)
	Upsert(ctx context.Context, userID uuid.UUID, taskType model.TaskType, date string, completed bool) (*model.Task, error)
}

func NewTaskService(taskRepo TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) GetToday(ctx context.Context, userID uuid.UUID) ([]*model.Task, error) {
	date := time.Now().UTC().Format("2006-01-02")
	tasks, err := s.taskRepo.GetToday(ctx, userID, date)
	if err != nil {
		return nil, err
	}
	taskMap := make(map[model.TaskType]*model.Task)
	for _, t := range tasks {
		taskMap[t.Type] = t
	}
	types := []model.TaskType{model.TaskCoding, model.TaskLeetcode, model.TaskContent}
	var out []*model.Task
	for _, typ := range types {
		if t, ok := taskMap[typ]; ok {
			out = append(out, t)
		} else {
			t, err := s.taskRepo.Upsert(ctx, userID, typ, date, false)
			if err != nil {
				return nil, err
			}
			out = append(out, t)
		}
	}
	return out, nil
}

func (s *TaskService) Update(ctx context.Context, userID uuid.UUID, taskType model.TaskType, date string, completed bool) (*model.Task, error) {
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	return s.taskRepo.Upsert(ctx, userID, taskType, date, completed)
}
