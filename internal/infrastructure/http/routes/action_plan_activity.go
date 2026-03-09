package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type ActionPlanActivityRoutes struct {
	activityHandler *api.ActionPlanActivityHandler
}

func NewActionPlanActivityRoutes(activityHandler *api.ActionPlanActivityHandler) *ActionPlanActivityRoutes {
	return &ActionPlanActivityRoutes{
		activityHandler: activityHandler,
	}
}

func (r *ActionPlanActivityRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	// Rotas diretas para activities (sem conflito)
	activities := v1.Group("/action-plan-activities")
	activities.Use(middleware.PartnerMiddleware())
	{
		activities.GET("/:id", r.activityHandler.GetByID)
		activities.PUT("/:id", r.activityHandler.Update)
		activities.DELETE("/:id", r.activityHandler.Delete)
		activities.POST("/:id/complete", r.activityHandler.MarkAsCompleted)
	}
}
