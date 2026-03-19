package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type TaskRepository interface {
	GetByDate(ctx context.Context, userID uuid.UUID, date string) ([]*model.Task, error)
	GetByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*model.Task, error)
	Create(ctx context.Context, userID uuid.UUID, title string, date string) (*model.Task, error)
	UpdateByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, title string, date string, completed bool) (*model.Task, error)
	DeleteByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) GetByDate(ctx context.Context, userID uuid.UUID, date string) ([]*model.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, "type", date, completed, created_at FROM tasks
		 WHERE user_id = $1 AND date = $2 ORDER BY created_at ASC`,
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

func (r *taskRepo) GetByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*model.Task, error) {
	t := &model.Task{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, "type", date, completed, created_at FROM tasks WHERE id = $1 AND user_id = $2`,
		taskID, userID,
	).Scan(&t.ID, &t.UserID, &t.Type, &t.Date, &t.Completed, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (r *taskRepo) Create(ctx context.Context, userID uuid.UUID, title string, date string) (*model.Task, error) {
	t := &model.Task{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tasks (user_id, "type", date, completed)
		 VALUES ($1, $2, $3, false)
		 RETURNING id, user_id, "type", date, completed, created_at`,
		userID, title, date,
	).Scan(&t.ID, &t.UserID, &t.Type, &t.Date, &t.Completed, &t.CreatedAt)
	return t, err
}

func (r *taskRepo) UpdateByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, title string, date string, completed bool) (*model.Task, error) {
	t := &model.Task{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE tasks SET "type" = $2, date = $3, completed = $4 WHERE id = $1 AND user_id = $5
		 RETURNING id, user_id, "type", date, completed, created_at`,
		taskID, title, date, completed, userID,
	).Scan(&t.ID, &t.UserID, &t.Type, &t.Date, &t.Completed, &t.CreatedAt)
	return t, err
}

func (r *taskRepo) DeleteByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1 AND user_id = $2`, taskID, userID)
	return err
}
