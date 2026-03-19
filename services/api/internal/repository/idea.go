package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"
	"strconv"

	"github.com/google/uuid"
)

type IdeaRepository interface {
	Create(ctx context.Context, userID uuid.UUID, hook, idea string, ideaType model.IdeaType, status model.IdeaStatus) (*model.Idea, error)
	List(ctx context.Context, userID uuid.UUID, filterType *model.IdeaType, filterStatus *model.IdeaStatus) ([]*model.Idea, error)
	Update(ctx context.Context, id, userID uuid.UUID, hook, idea *string, ideaType *model.IdeaType, status *model.IdeaStatus) (*model.Idea, error)
}

type ideaRepo struct {
	db *sql.DB
}

func NewIdeaRepository(db *sql.DB) IdeaRepository {
	return &ideaRepo{db: db}
}

func (r *ideaRepo) Create(ctx context.Context, userID uuid.UUID, hook, idea string, ideaType model.IdeaType, status model.IdeaStatus) (*model.Idea, error) {
	row := &model.Idea{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO ideas (user_id, hook, idea, "type", status)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, hook, idea, "type", status, created_at`,
		userID, hook, idea, ideaType, status,
	).Scan(&row.ID, &row.UserID, &row.Hook, &row.Idea, &row.Type, &row.Status, &row.CreatedAt)
	return row, err
}

func (r *ideaRepo) List(ctx context.Context, userID uuid.UUID, filterType *model.IdeaType, filterStatus *model.IdeaStatus) ([]*model.Idea, error) {
	query := `SELECT id, user_id, hook, idea, "type", status, created_at FROM ideas WHERE user_id = $1`
	args := []interface{}{userID}
	pos := 2
	if filterType != nil {
		query += ` AND "type" = $` + strconv.Itoa(pos)
		args = append(args, *filterType)
		pos++
	}
	if filterStatus != nil {
		query += ` AND status = $` + strconv.Itoa(pos)
		args = append(args, *filterStatus)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Idea
	for rows.Next() {
		i := &model.Idea{}
		if err := rows.Scan(&i.ID, &i.UserID, &i.Hook, &i.Idea, &i.Type, &i.Status, &i.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

func (r *ideaRepo) Update(ctx context.Context, id, userID uuid.UUID, hook, idea *string, ideaType *model.IdeaType, status *model.IdeaStatus) (*model.Idea, error) {
	row := &model.Idea{}
	var typeVal, statusVal interface{}
	if ideaType != nil {
		typeVal = string(*ideaType)
	}
	if status != nil {
		statusVal = string(*status)
	}
	err := r.db.QueryRowContext(ctx,
		`UPDATE ideas SET
			hook = COALESCE($2, hook),
			idea = COALESCE($3, idea),
			"type" = COALESCE($4, "type"),
			status = COALESCE($5, status)
		 WHERE id = $1 AND user_id = $6
		 RETURNING id, user_id, hook, idea, "type", status, created_at`,
		id, hook, idea, typeVal, statusVal, userID,
	).Scan(&row.ID, &row.UserID, &row.Hook, &row.Idea, &row.Type, &row.Status, &row.CreatedAt)
	return row, err
}
