package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"
	"fmt"

	"github.com/google/uuid"
)

type FinanceRepository interface {
	Create(ctx context.Context, userID uuid.UUID, amount float64, financeType model.FinanceType, note, date string) (*model.Finance, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.Finance, error)
	SumByMonth(ctx context.Context, userID uuid.UUID, year int, month int) (float64, error)
}

type financeRepo struct {
	db *sql.DB
}

func NewFinanceRepository(db *sql.DB) FinanceRepository {
	return &financeRepo{db: db}
}

func (r *financeRepo) Create(ctx context.Context, userID uuid.UUID, amount float64, financeType model.FinanceType, note, date string) (*model.Finance, error) {
	f := &model.Finance{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO finances (user_id, amount, "type", note, date)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, amount, "type", note, date, created_at`,
		userID, amount, financeType, note, date,
	).Scan(&f.ID, &f.UserID, &f.Amount, &f.Type, &f.Note, &f.Date, &f.CreatedAt)
	return f, err
}

func (r *financeRepo) List(ctx context.Context, userID uuid.UUID) ([]*model.Finance, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, amount, "type", note, date, created_at
		 FROM finances WHERE user_id = $1 ORDER BY date DESC, created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Finance
	for rows.Next() {
		f := &model.Finance{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.Amount, &f.Type, &f.Note, &f.Date, &f.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *financeRepo) SumByMonth(ctx context.Context, userID uuid.UUID, year, month int) (float64, error) {
	var sum sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM finances
		 WHERE user_id = $1 AND date >= $2::date AND date < $2::date + interval '1 month'`,
		userID, formatYearMonth(year, month),
	).Scan(&sum)
	if err != nil || !sum.Valid {
		return 0, err
	}
	return sum.Float64, nil
}

func formatYearMonth(year, month int) string {
	return fmt.Sprintf("%04d-%02d-01", year, month)
}
