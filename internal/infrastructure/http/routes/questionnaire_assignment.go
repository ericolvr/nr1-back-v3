package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type QuestionnaireAssignmentRoutes struct {
	assignmentHandler *api.QuestionnaireAssignmentHandler
}

func NewQuestionnaireAssignmentRoutes(assignmentHandler *api.QuestionnaireAssignmentHandler) *QuestionnaireAssignmentRoutes {
	return &QuestionnaireAssignmentRoutes{
		assignmentHandler: assignmentHandler,
	}
}

func (r *QuestionnaireAssignmentRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	assignments := v1.Group("/questionnaire-assignments")
	assignments.Use(middleware.PartnerMiddleware())
	{
		assignments.POST("", r.assignmentHandler.Create)
		assignments.GET("", r.assignmentHandler.List)
		assignments.GET("/:id", r.assignmentHandler.GetByID)
		assignments.POST("/:id/close", r.assignmentHandler.Close)
		assignments.DELETE("/:id", r.assignmentHandler.Delete)
	}
}
