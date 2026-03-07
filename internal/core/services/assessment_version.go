package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentVersionService struct {
	versionRepo domain.AssessmentVersionRepository
}

func NewAssessmentVersionService(versionRepo domain.AssessmentVersionRepository) *AssessmentVersionService {
	return &AssessmentVersionService{
		versionRepo: versionRepo,
	}
}

type TemplateChanges struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

func (s *AssessmentVersionService) CreateVersion(ctx context.Context, templateID, partnerID int64, version int, changes []TemplateChanges, createdBy int64) error {
	if templateID <= 0 || partnerID <= 0 {
		return errors.New("invalid IDs")
	}

	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %w", err)
	}

	assessmentVersion := &domain.AssessmentVersion{
		TemplateID: templateID,
		PartnerID:  partnerID,
		Version:    version,
		Changes:    string(changesJSON),
		CreatedBy:  createdBy,
	}

	if err := assessmentVersion.Validate(); err != nil {
		return err
	}

	return s.versionRepo.Create(ctx, assessmentVersion)
}

func (s *AssessmentVersionService) ListByTemplate(ctx context.Context, partnerID, templateID int64, limit, offset int64) ([]*domain.AssessmentVersion, error) {
	if partnerID <= 0 || templateID <= 0 {
		return nil, errors.New("invalid IDs")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.versionRepo.ListByTemplate(ctx, partnerID, templateID, limit, offset)
}

func (s *AssessmentVersionService) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentVersion, error) {
	if partnerID <= 0 || id <= 0 {
		return nil, errors.New("invalid IDs")
	}

	return s.versionRepo.GetByID(ctx, partnerID, id)
}

func (s *AssessmentVersionService) TrackTemplateUpdate(ctx context.Context, oldTemplate, newTemplate *domain.AssessmentTemplate, userID int64) error {
	var changes []TemplateChanges

	fmt.Printf("[AUDIT] Comparing templates:\n")
	fmt.Printf("  Old: Name=%s, Desc=%s, Active=%v\n", oldTemplate.Name, oldTemplate.Description, oldTemplate.Active)
	fmt.Printf("  New: Name=%s, Desc=%s, Active=%v\n", newTemplate.Name, newTemplate.Description, newTemplate.Active)

	if oldTemplate.Name != newTemplate.Name {
		changes = append(changes, TemplateChanges{
			Field:    "name",
			OldValue: oldTemplate.Name,
			NewValue: newTemplate.Name,
		})
	}

	if oldTemplate.Description != newTemplate.Description {
		changes = append(changes, TemplateChanges{
			Field:    "description",
			OldValue: oldTemplate.Description,
			NewValue: newTemplate.Description,
		})
	}

	if oldTemplate.Active != newTemplate.Active {
		changes = append(changes, TemplateChanges{
			Field:    "active",
			OldValue: oldTemplate.Active,
			NewValue: newTemplate.Active,
		})
	}

	fmt.Printf("[AUDIT] Detected %d changes\n", len(changes))
	if len(changes) == 0 {
		fmt.Printf("[AUDIT] No changes detected, skipping version creation\n")
		return nil
	}

	fmt.Printf("[AUDIT] Creating version record...\n")
	return s.CreateVersion(ctx, newTemplate.ID, newTemplate.PartnerID, newTemplate.Version, changes, userID)
}
