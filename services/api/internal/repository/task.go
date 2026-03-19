package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type TaskRepository interface {
	GetToday(ctx context.Context, userID uuid.UUID, date string) ([]*model.Task, error)
	Upsert(ctx context.Context, userID uuid.UUID, taskType model.TaskType, date string, completed bool) (*model.Task, error)
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) GetToday(ctx context.Context, userID uuid.UUID, date string) ([]*model.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, "type", date, completed, created_at FROM tasks
		 WHERE user_id = $1 AND date = $2 ORDER BY "type"`,
		userID, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Task
	for rows.Next() {
		t := &model.Task{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Type, &t.Date, &t.Completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *taskRepo) Upsert(ctx context.Context, userID uuid.UUID, taskType model.TaskType, date string, completed bool) (*model.Task, error) {
	t := &model.Task{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tasks (user_id, "type", date, completed)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id, "type", date) DO UPDATE SET completed = $4
		 RETURNING id, user_id, "type", date, completed, created_at`,
		userID, taskType, date, completed,
	).Scan(&t.ID, &t.UserID, &t.Type, &t.Date, &t.Completed, &t.CreatedAt)
	return t, err
}
