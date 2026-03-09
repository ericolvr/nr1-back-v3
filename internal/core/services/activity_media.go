package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type ActivityMediaService struct {
	mediaRepo domain.ActivityMediaRepository
}

func NewActivityMediaService(mediaRepo domain.ActivityMediaRepository) *ActivityMediaService {
	return &ActivityMediaService{
		mediaRepo: mediaRepo,
	}
}

func (s *ActivityMediaService) Create(ctx context.Context, media *domain.ActivityMedia) error {
	if err := media.Validate(); err != nil {
		return err
	}
	return s.mediaRepo.Create(ctx, media)
}

func (s *ActivityMediaService) GetByID(ctx context.Context, id string) (*domain.ActivityMedia, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("id must be a valid number")
	}

	return s.mediaRepo.GetByID(ctx, idInt)
}

func (s *ActivityMediaService) GetByActivityID(ctx context.Context, activityID int64) ([]*domain.ActivityMedia, error) {
	return s.mediaRepo.GetByActivityID(ctx, activityID)
}

func (s *ActivityMediaService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("id must be a valid number")
	}

	return s.mediaRepo.Delete(ctx, idInt)
}
