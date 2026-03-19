package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, userID uuid.UUID, startTime time.Time) (*model.Session, error)
	GetActive(ctx context.Context, userID uuid.UUID) (*model.Session, error)
	End(ctx context.Context, id, userID uuid.UUID, endTime time.Time) (*model.Session, error)
}

type sessionRepo struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) Create(ctx context.Context, userID uuid.UUID, startTime time.Time) (*model.Session, error) {
	s := &model.Session{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO sessions (user_id, start_time) VALUES ($1, $2)
		 RETURNING id, user_id, start_time, end_time`,
		userID, startTime,
	).Scan(&s.ID, &s.UserID, &s.StartTime, &s.EndTime)
	return s, err
}

func (r *sessionRepo) GetActive(ctx context.Context, userID uuid.UUID) (*model.Session, error) {
	s := &model.Session{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, start_time, end_time FROM sessions
		 WHERE user_id = $1 AND end_time IS NULL ORDER BY start_time DESC LIMIT 1`,
		userID,
	).Scan(&s.ID, &s.UserID, &s.StartTime, &s.EndTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *sessionRepo) End(ctx context.Context, id, userID uuid.UUID, endTime time.Time) (*model.Session, error) {
	s := &model.Session{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE sessions SET end_time = $3 WHERE id = $1 AND user_id = $2 AND end_time IS NULL
		 RETURNING id, user_id, start_time, end_time`,
		id, userID, endTime,
	).Scan(&s.ID, &s.UserID, &s.StartTime, &s.EndTime)
	return s, err
}
