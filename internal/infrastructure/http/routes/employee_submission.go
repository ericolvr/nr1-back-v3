package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type EmployeeSubmissionRoutes struct {
	submissionHandler *api.EmployeeSubmissionHandler
}

func NewEmployeeSubmissionRoutes(submissionHandler *api.EmployeeSubmissionHandler) *EmployeeSubmissionRoutes {
	return &EmployeeSubmissionRoutes{
		submissionHandler: submissionHandler,
	}
}

func (r *EmployeeSubmissionRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	// Public routes (no middleware)
	responses := v1.Group("/responses")
	{
		responses.GET("/:token", r.submissionHandler.GetByToken)
		responses.POST("/:token/submit", r.submissionHandler.Submit)
	}

	// Protected routes
	submissions := v1.Group("/submissions")
	submissions.Use(middleware.PartnerMiddleware())
	{
		submissions.GET("", r.submissionHandler.List)
	}
}
