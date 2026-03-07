package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentTemplateService struct {
	templateRepo   domain.AssessmentTemplateRepository
	partnerRepo    domain.PartnerRepository
	versionService *AssessmentVersionService
}

func NewAssessmentTemplateService(
	templateRepo domain.AssessmentTemplateRepository,
	partnerRepo domain.PartnerRepository,
	versionService *AssessmentVersionService,
) *AssessmentTemplateService {
	return &AssessmentTemplateService{
		templateRepo:   templateRepo,
		partnerRepo:    partnerRepo,
		versionService: versionService,
	}
}

func (s *AssessmentTemplateService) Create(ctx context.Context, template *domain.AssessmentTemplate) error {
	if err := domain.ValidateAssessmentTemplate(template); err != nil {
		return err
	}

	partner, err := s.partnerRepo.GetByID(ctx, template.PartnerID)
	if err != nil {
		return errors.New("partner not found")
	}
	if !partner.Active {
		return errors.New("partner is not active")
	}

	if template.Version == 0 {
		template.Version = 1
	}

	return s.templateRepo.Create(ctx, template)
}

func (s *AssessmentTemplateService) List(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.AssessmentTemplate, error) {
	if partnerID <= 0 {
		return nil, errors.New("invalid partner ID")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.templateRepo.List(ctx, partnerID, limit, offset)
}

func (s *AssessmentTemplateService) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentTemplate, error) {
	if partnerID <= 0 || id <= 0 {
		return nil, errors.New("invalid IDs")
	}

	return s.templateRepo.GetByID(ctx, partnerID, id)
}

func (s *AssessmentTemplateService) Update(ctx context.Context, template *domain.AssessmentTemplate, userID int64) error {
	if template.ID <= 0 || template.PartnerID <= 0 {
		return errors.New("invalid IDs")
	}

	if err := domain.ValidateAssessmentTemplate(template); err != nil {
		return err
	}

	oldTemplate, err := s.templateRepo.GetByID(ctx, template.PartnerID, template.ID)
	if err != nil {
		return err
	}

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return err
	}

	fmt.Printf("[AUDIT] versionService=%v, userID=%d\n", s.versionService != nil, userID)
	if s.versionService != nil && userID > 0 {
		fmt.Printf("[AUDIT] Tracking update for template ID=%d, PartnerID=%d, UserID=%d\n", template.ID, template.PartnerID, userID)
		if err := s.versionService.TrackTemplateUpdate(ctx, oldTemplate, template, userID); err != nil {
			fmt.Printf("[AUDIT ERROR] Failed to track: %v\n", err)
			return err
		}
		fmt.Printf("[AUDIT] Successfully tracked update\n")
	} else if s.versionService != nil {
		fmt.Printf("[AUDIT] Skipping version tracking - userID is 0 (no authentication)\n")
	}

	return nil
}

func (s *AssessmentTemplateService) IncrementVersion(ctx context.Context, partnerID, id int64) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.templateRepo.IncrementVersion(ctx, partnerID, id)
}

func (s *AssessmentTemplateService) Delete(ctx context.Context, partnerID, id int64) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.templateRepo.Delete(ctx, partnerID, id)
}

func (s *AssessmentTemplateService) ToggleActive(ctx context.Context, partnerID, id int64, active bool) error {
	if partnerID <= 0 || id <= 0 {
		return errors.New("invalid IDs")
	}

	return s.templateRepo.ToggleActive(ctx, partnerID, id, active)
}

func (s *AssessmentTemplateService) ListDeleted(ctx context.Context, partnerID int64, limit, offset int64) ([]*domain.AssessmentTemplate, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.templateRepo.ListDeleted(ctx, partnerID, limit, offset)
}
