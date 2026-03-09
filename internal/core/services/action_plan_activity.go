package services

import (
	"context"
	"fmt"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type ActionPlanActivityService struct {
	activityRepo domain.ActionPlanActivityRepository
	actionPlanRepo domain.ActionPlanRepository
}

func NewActionPlanActivityService(
	activityRepo domain.ActionPlanActivityRepository,
	actionPlanRepo domain.ActionPlanRepository,
) *ActionPlanActivityService {
	return &ActionPlanActivityService{
		activityRepo: activityRepo,
		actionPlanRepo: actionPlanRepo,
	}
}

func (s *ActionPlanActivityService) Create(ctx context.Context, activity *domain.ActionPlanActivity) error {
	// Validar activity
	if err := activity.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Verificar se o action plan existe
	_, err := s.actionPlanRepo.GetByID(ctx, 0, activity.ActionPlanID) // TODO: passar partnerID correto
	if err != nil {
		return fmt.Errorf("action plan not found: %w", err)
	}

	// Criar activity
	if err := s.activityRepo.Create(ctx, activity); err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	return nil
}

func (s *ActionPlanActivityService) GetByID(ctx context.Context, id int64) (*domain.ActionPlanActivity, error) {
	activity, err := s.activityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}
	return activity, nil
}

func (s *ActionPlanActivityService) ListByActionPlan(ctx context.Context, actionPlanID int64) ([]*domain.ActionPlanActivity, error) {
	activities, err := s.activityRepo.ListByActionPlan(ctx, actionPlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to list activities: %w", err)
	}
	return activities, nil
}

func (s *ActionPlanActivityService) Update(ctx context.Context, activity *domain.ActionPlanActivity) error {
	// Validar activity
	if err := activity.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Verificar se existe
	_, err := s.activityRepo.GetByID(ctx, activity.ID)
	if err != nil {
		return fmt.Errorf("activity not found: %w", err)
	}

	// Atualizar
	if err := s.activityRepo.Update(ctx, activity); err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	return nil
}

func (s *ActionPlanActivityService) Delete(ctx context.Context, id int64) error {
	// Verificar se existe
	_, err := s.activityRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("activity not found: %w", err)
	}

	// Deletar
	if err := s.activityRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	return nil
}

func (s *ActionPlanActivityService) MarkAsCompleted(ctx context.Context, id int64) error {
	// Buscar activity
	activity, err := s.activityRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("activity not found: %w", err)
	}

	// Marcar como concluída
	activity.MarkAsCompleted()

	// Atualizar
	if err := s.activityRepo.Update(ctx, activity); err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	return nil
}
