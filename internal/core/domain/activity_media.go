package domain

import (
	"context"
	"errors"
	"time"
)

type ActivityMedia struct {
	ID         int64     `json:"id"`
	ActivityID int64     `json:"activity_id"`
	MediaURL   string    `json:"media_url"`
	MediaType  string    `json:"media_type"`
	CreatedAt  time.Time `json:"created_at"`
}

type ActivityMediaRepository interface {
	Create(ctx context.Context, media *ActivityMedia) error
	GetByID(ctx context.Context, id int64) (*ActivityMedia, error)
	GetByActivityID(ctx context.Context, activityID int64) ([]*ActivityMedia, error)
	Delete(ctx context.Context, id int64) error
}

func (am *ActivityMedia) Validate() error {
	if am.ActivityID <= 0 {
		return errors.New("activity_id is required")
	}

	if am.MediaURL == "" {
		return errors.New("media_url is required")
	}

	if am.MediaType == "" {
		return errors.New("media_type is required")
	}

	validMediaTypes := map[string]bool{
		"photo":    true,
		"video":    true,
		"document": true,
	}
	if !validMediaTypes[am.MediaType] {
		return errors.New("media_type must be one of: photo, video, document")
	}

	return nil
}
