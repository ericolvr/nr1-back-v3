package api

import (
	"net/http"

	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type SurveyHandler struct {
	submissionService *services.EmployeeSubmissionService
	questionService   *services.QuestionService
	templateService   *services.AssessmentTemplateService
}

func NewSurveyHandler(
	submissionService *services.EmployeeSubmissionService,
	questionService *services.QuestionService,
	templateService *services.AssessmentTemplateService,
) *SurveyHandler {
	return &SurveyHandler{
		submissionService: submissionService,
		questionService:   questionService,
		templateService:   templateService,
	}
}

func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token é obrigatório",
		})
		return
	}

	// Buscar EmployeeSubmission pelo token
	submission, err := h.submissionService.GetByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Token inválido ou expirado",
		})
		return
	}

	// Buscar AssessmentTemplate
	template, err := h.templateService.GetByID(c.Request.Context(), submission.PartnerID, submission.TemplateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar template",
		})
		return
	}

	// Buscar Questions do template
	questions, err := h.questionService.List(c.Request.Context(), submission.PartnerID, submission.TemplateID, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar perguntas",
		})
		return
	}

	// Montar resposta
	var questionDTOs []dto.QuestionResponse
	for _, q := range questions {
		questionDTOs = append(questionDTOs, dto.QuestionResponse{
			ID:              q.ID,
			PartnerID:       q.PartnerID,
			TemplateID:      q.TemplateID,
			Question:        q.Question,
			Type:            q.Type,
			Options:         q.Options,
			ScoreValues:     q.ScoreValues,
			Weight:          q.Weight,
			Required:        q.Required,
			OrderNum:        q.OrderNum,
			Category:        q.Category,
			CreatedAt:       q.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       q.UpdatedAt.Format("2006-01-02 15:04:05"),
			TemplateName:    template.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"submission_id":         submission.ID,
		"template_id":           submission.TemplateID,
		"template_name":         template.Name,
		"template_description":  template.Description,
		"department_id":         submission.DepartmentID,
		"status":                submission.Status,
		"questions":             questionDTOs,
	})
}
