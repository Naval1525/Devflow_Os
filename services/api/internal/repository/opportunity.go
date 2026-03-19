package repository

import (
	"context"
	"database/sql"
	"devflowos/api/internal/model"

	"github.com/google/uuid"
)

type OpportunityRepository interface {
	Create(ctx context.Context, userID uuid.UUID, name string, oppType model.OpportunityType, stage model.OpportunityStage, source, notes string) (*model.Opportunity, error)
	List(ctx context.Context, userID uuid.UUID) ([]*model.Opportunity, error)
	Update(ctx context.Context, id, userID uuid.UUID, name, source, notes *string, oppType *model.OpportunityType, stage *model.OpportunityStage) (*model.Opportunity, error)
}

type opportunityRepo struct {
	db *sql.DB
}

func NewOpportunityRepository(db *sql.DB) OpportunityRepository {
	return &opportunityRepo{db: db}
}

func (r *opportunityRepo) Create(ctx context.Context, userID uuid.UUID, name string, oppType model.OpportunityType, stage model.OpportunityStage, source, notes string) (*model.Opportunity, error) {
	o := &model.Opportunity{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO opportunities (user_id, name, "type", stage, source, notes)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, name, "type", stage, source, notes, created_at`,
		userID, name, oppType, stage, source, notes,
	).Scan(&o.ID, &o.UserID, &o.Name, &o.Type, &o.Stage, &o.Source, &o.Notes, &o.CreatedAt)
	return o, err
}

func (r *opportunityRepo) List(ctx context.Context, userID uuid.UUID) ([]*model.Opportunity, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, name, "type", stage, source, notes, created_at
		 FROM opportunities WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Opportunity
	for rows.Next() {
		o := &model.Opportunity{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.Name, &o.Type, &o.Stage, &o.Source, &o.Notes, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func (r *opportunityRepo) Update(ctx context.Context, id, userID uuid.UUID, name, source, notes *string, oppType *model.OpportunityType, stage *model.OpportunityStage) (*model.Opportunity, error) {
	o := &model.Opportunity{}
	var typeVal, stageVal interface{}
	if oppType != nil {
		typeVal = string(*oppType)
	}
	if stage != nil {
		stageVal = string(*stage)
	}
	err := r.db.QueryRowContext(ctx,
		`UPDATE opportunities SET
			name = COALESCE($2, name),
			source = COALESCE($3, source),
			notes = COALESCE($4, notes),
			"type" = COALESCE($5, "type"),
			stage = COALESCE($6, stage)
		 WHERE id = $1 AND user_id = $7
		 RETURNING id, user_id, name, "type", stage, source, notes, created_at`,
		id, name, source, notes, typeVal, stageVal, userID,
	).Scan(&o.ID, &o.UserID, &o.Name, &o.Type, &o.Stage, &o.Source, &o.Notes, &o.CreatedAt)
	return o, err
}
