package database

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/lib/pq"
)

type ActionPlanActivityRepository struct {
	db *sql.DB
}

func NewActionPlanActivityRepository(db *sql.DB) *ActionPlanActivityRepository {
	return &ActionPlanActivityRepository{db: db}
}

func (r *ActionPlanActivityRepository) Create(ctx context.Context, activity *domain.ActionPlanActivity) error {
	query := `
		INSERT INTO action_plan_activities (
			action_plan_id, type, title, description, status, 
			due_date, completed_at, evidence_urls, created_by, created_by_name
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	evidenceURLsJSON, _ := json.Marshal(activity.EvidenceURLs)

	return r.db.QueryRowContext(
		ctx, query,
		activity.ActionPlanID,
		activity.Type,
		activity.Title,
		activity.Description,
		activity.Status,
		activity.DueDate,
		activity.CompletedAt,
		evidenceURLsJSON,
		activity.CreatedBy,
		activity.CreatedByName,
	).Scan(&activity.ID, &activity.CreatedAt, &activity.UpdatedAt)
}

func (r *ActionPlanActivityRepository) GetByID(ctx context.Context, id int64) (*domain.ActionPlanActivity, error) {
	query := `
		SELECT id, action_plan_id, type, title, description, status,
			due_date, completed_at, evidence_urls, created_by, created_by_name,
			created_at, updated_at
		FROM action_plan_activities
		WHERE id = $1
	`

	activity := &domain.ActionPlanActivity{}
	var evidenceURLsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.ActionPlanID,
		&activity.Type,
		&activity.Title,
		&activity.Description,
		&activity.Status,
		&activity.DueDate,
		&activity.CompletedAt,
		&evidenceURLsJSON,
		&activity.CreatedBy,
		&activity.CreatedByName,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if len(evidenceURLsJSON) > 0 {
		json.Unmarshal(evidenceURLsJSON, &activity.EvidenceURLs)
	}

	return activity, nil
}

func (r *ActionPlanActivityRepository) ListByActionPlan(ctx context.Context, actionPlanID int64) ([]*domain.ActionPlanActivity, error) {
	query := `
		SELECT id, action_plan_id, type, title, description, status,
			due_date, completed_at, evidence_urls, created_by, created_by_name,
			created_at, updated_at
		FROM action_plan_activities
		WHERE action_plan_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, actionPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*domain.ActionPlanActivity

	for rows.Next() {
		activity := &domain.ActionPlanActivity{}
		var evidenceURLsJSON []byte

		err := rows.Scan(
			&activity.ID,
			&activity.ActionPlanID,
			&activity.Type,
			&activity.Title,
			&activity.Description,
			&activity.Status,
			&activity.DueDate,
			&activity.CompletedAt,
			&evidenceURLsJSON,
			&activity.CreatedBy,
			&activity.CreatedByName,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if len(evidenceURLsJSON) > 0 {
			json.Unmarshal(evidenceURLsJSON, &activity.EvidenceURLs)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *ActionPlanActivityRepository) Update(ctx context.Context, activity *domain.ActionPlanActivity) error {
	query := `
		UPDATE action_plan_activities
		SET title = $1, description = $2, status = $3, due_date = $4,
			completed_at = $5, evidence_urls = $6, updated_at = NOW()
		WHERE id = $7
	`

	evidenceURLsJSON, _ := json.Marshal(activity.EvidenceURLs)

	_, err := r.db.ExecContext(
		ctx, query,
		activity.Title,
		activity.Description,
		activity.Status,
		activity.DueDate,
		activity.CompletedAt,
		evidenceURLsJSON,
		activity.ID,
	)

	return err
}

func (r *ActionPlanActivityRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM action_plan_activities WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// scanActivity é um helper para escanear uma activity do banco
func scanActivity(row *sql.Row) (*domain.ActionPlanActivity, error) {
	activity := &domain.ActionPlanActivity{}
	var evidenceURLs pq.StringArray

	err := row.Scan(
		&activity.ID,
		&activity.ActionPlanID,
		&activity.Type,
		&activity.Title,
		&activity.Description,
		&activity.Status,
		&activity.DueDate,
		&activity.CompletedAt,
		&evidenceURLs,
		&activity.CreatedBy,
		&activity.CreatedByName,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	activity.EvidenceURLs = evidenceURLs
	return activity, nil
}
