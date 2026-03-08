package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
)

type AssessmentAssignmentService struct {
	assignmentRepo domain.AssessmentAssignmentRepository
	departmentRepo domain.DepartmentRepository
}

func NewAssessmentAssignmentService(
	assignmentRepo domain.AssessmentAssignmentRepository,
	departmentRepo domain.DepartmentRepository,
) *AssessmentAssignmentService {
	return &AssessmentAssignmentService{
		assignmentRepo: assignmentRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *AssessmentAssignmentService) Create(ctx context.Context, assignment *domain.AssessmentAssignment) error {
	if err := assignment.Validate(); err != nil {
		return err
	}

	// Validar que todos os departments existem
	for _, deptID := range assignment.DepartmentIDs {
		_, err := s.departmentRepo.GetByID(ctx, assignment.PartnerID, deptID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("department not found")
			}
			return err
		}
	}

	return s.assignmentRepo.Create(ctx, assignment)
}

func (s *AssessmentAssignmentService) GetByID(ctx context.Context, partnerID, id int64) (*domain.AssessmentAssignment, error) {
	return s.assignmentRepo.GetByID(ctx, partnerID, id)
}

func (s *AssessmentAssignmentService) GetByTemplateAndDepartment(ctx context.Context, partnerID, templateID, departmentID int64) (*domain.AssessmentAssignment, error) {
	return s.assignmentRepo.GetByTemplateAndDepartment(ctx, partnerID, templateID, departmentID)
}

func (s *AssessmentAssignmentService) List(ctx context.Context, partnerID, limit, offset int64) ([]*domain.AssessmentAssignment, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.assignmentRepo.List(ctx, partnerID, limit, offset)
}

func (s *AssessmentAssignmentService) ListByTemplate(ctx context.Context, partnerID, templateID int64) ([]*domain.AssessmentAssignment, error) {
	return s.assignmentRepo.ListByTemplate(ctx, partnerID, templateID, MaxAssignmentsPerQuestionnaire, 0)
}

func (s *AssessmentAssignmentService) Deactivate(ctx context.Context, partnerID, id int64) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	assignment.Active = false
	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *AssessmentAssignmentService) Activate(ctx context.Context, partnerID, id int64) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return err
	}

	assignment.Active = true
	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *AssessmentAssignmentService) Update(ctx context.Context, partnerID, id int64, templateID *int64, departmentIDs []int64, active *bool) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, partnerID, id)
	if err != nil {
		return errors.New("assignment not found")
	}

	// Atualizar campos se fornecidos
	if templateID != nil {
		assignment.TemplateID = *templateID
	}
	if len(departmentIDs) > 0 {
		// Validar que todos os departments existem
		for _, deptID := range departmentIDs {
			_, err := s.departmentRepo.GetByID(ctx, partnerID, deptID)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("department not found")
				}
				return err
			}
		}
		assignment.DepartmentIDs = departmentIDs
	}
	if active != nil {
		assignment.Active = *active
	}

	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *AssessmentAssignmentService) Delete(ctx context.Context, partnerID, id int64) error {
	return s.assignmentRepo.Delete(ctx, partnerID, id)
}
