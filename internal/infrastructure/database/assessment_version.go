package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentVersionRepository struct {
	db *sql.DB
}

func NewAssessmentVersionRepository(db *sql.DB) *AssessmentVersionRepository {
	return &AssessmentVersionRepository{db: db}
}

func (r *AssessmentVersionRepository) Create(ctx context.Context, version *domain.AssessmentVersion) error {
	query := `
		INSERT INTO assessment_versions (template_id, partner_id, version, changes, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		version.TemplateID,
		version.PartnerID,
		version.Version,
		version.Changes,
		version.CreatedBy,
	).Scan(&version.ID, &version.CreatedAt)
	return err
}

func (r *AssessmentVersionRepository) ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*domain.AssessmentVersion, error) {
	query := `
		SELECT id, template_id, partner_id, version, changes, created_by, created_at
		FROM assessment_versions
		WHERE partner_id = $1 AND template_id = $2
		ORDER BY version DESC, created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, query, partnerID, templateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*domain.AssessmentVersion
	for rows.Next() {
		var v domain.AssessmentVersion
		var createdBy sql.NullInt64
		err := rows.Scan(
			&v.ID,
			&v.TemplateID,
			&v.PartnerID,
			&v.Version,
			&v.Changes,
			&createdBy,
			&v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if createdBy.Valid {
			v.CreatedBy = createdBy.Int64
		}
		versions = append(versions, &v)
	}
	return versions, nil
}

func (r *AssessmentVersionRepository) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentVersion, error) {
	query := `
		SELECT id, template_id, partner_id, version, changes, created_by, created_at
		FROM assessment_versions
		WHERE partner_id = $1 AND id = $2
	`
	var v domain.AssessmentVersion
	var createdBy sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, partnerID, id).Scan(
		&v.ID,
		&v.TemplateID,
		&v.PartnerID,
		&v.Version,
		&v.Changes,
		&createdBy,
		&v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if createdBy.Valid {
		v.CreatedBy = createdBy.Int64
	}
	return &v, nil
}
