package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type LeetCodeRepository interface {
	Create(ctx context.Context, userID uuid.UUID, problemName string, difficulty model.Difficulty, approach, mistake string, timeTaken *int) (*model.LeetCodeLog, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.LeetCodeLog, error)
}

type leetcodeRepo struct {
	db *sql.DB
}

func NewLeetCodeRepository(db *sql.DB) LeetCodeRepository {
	return &leetcodeRepo{db: db}
}

func (r *leetcodeRepo) Create(ctx context.Context, userID uuid.UUID, problemName string, difficulty model.Difficulty, approach, mistake string, timeTaken *int) (*model.LeetCodeLog, error) {
	row := &model.LeetCodeLog{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO leetcode_logs (user_id, problem_name, difficulty, approach, mistake, time_taken)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, problem_name, difficulty, approach, mistake, time_taken, created_at`,
		userID, problemName, difficulty, approach, mistake, timeTaken,
	).Scan(&row.ID, &row.UserID, &row.ProblemName, &row.Difficulty, &row.Approach, &row.Mistake, &row.TimeTaken, &row.CreatedAt)
	return row, err
}

func (r *leetcodeRepo) List(ctx context.Context, userID uuid.UUID) ([]*model.LeetCodeLog, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, problem_name, difficulty, approach, mistake, time_taken, created_at
		 FROM leetcode_logs WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.LeetCodeLog
	for rows.Next() {
		l := &model.LeetCodeLog{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.ProblemName, &l.Difficulty, &l.Approach, &l.Mistake, &l.TimeTaken, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
