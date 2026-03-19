package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type CodingLogRepository interface {
	Create(ctx context.Context, userID uuid.UUID, title, description string) (*model.CodingLog, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.CodingLog, error)
}

type codingLogRepo struct {
	db *sql.DB
}

func NewCodingLogRepository(db *sql.DB) CodingLogRepository {
	return &codingLogRepo{db: db}
}

func (r *codingLogRepo) Create(ctx context.Context, userID uuid.UUID, title, description string) (*model.CodingLog, error) {
	log := &model.CodingLog{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO coding_logs (user_id, title, description)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, title, description, created_at`,
		userID, title, description,
	).Scan(&log.ID, &log.UserID, &log.Title, &log.Description, &log.CreatedAt)
	return log, err
}

func (r *codingLogRepo) List(ctx context.Context, userID uuid.UUID) ([]*model.CodingLog, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, description, created_at
		 FROM coding_logs WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.CodingLog
	for rows.Next() {
		log := &model.CodingLog{}
		if err := rows.Scan(&log.ID, &log.UserID, &log.Title, &log.Description, &log.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, log)
	}
	return out, rows.Err()
}
