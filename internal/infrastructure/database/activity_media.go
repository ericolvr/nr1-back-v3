package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type ActivityMediaRepository struct {
	db *sql.DB
}

func NewActivityMediaRepository(db *sql.DB) *ActivityMediaRepository {
	return &ActivityMediaRepository{db: db}
}

func (r *ActivityMediaRepository) Create(ctx context.Context, media *domain.ActivityMedia) error {
	query := `
		INSERT INTO activity_media (activity_id, media_url, media_type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		media.ActivityID,
		media.MediaURL,
		media.MediaType,
	).Scan(&media.ID, &media.CreatedAt)
}

func (r *ActivityMediaRepository) GetByID(ctx context.Context, id int64) (*domain.ActivityMedia, error) {
	query := `
		SELECT id, activity_id, media_url, media_type, created_at
		FROM activity_media
		WHERE id = $1`

	media := &domain.ActivityMedia{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&media.ID,
		&media.ActivityID,
		&media.MediaURL,
		&media.MediaType,
		&media.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return media, nil
}

func (r *ActivityMediaRepository) GetByActivityID(ctx context.Context, activityID int64) ([]*domain.ActivityMedia, error) {
	query := `
		SELECT id, activity_id, media_url, media_type, created_at
		FROM activity_media
		WHERE activity_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []*domain.ActivityMedia
	for rows.Next() {
		media := &domain.ActivityMedia{}
		err := rows.Scan(
			&media.ID,
			&media.ActivityID,
			&media.MediaURL,
			&media.MediaType,
			&media.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		mediaList = append(mediaList, media)
	}
	return mediaList, nil
}

func (r *ActivityMediaRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM activity_media WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
