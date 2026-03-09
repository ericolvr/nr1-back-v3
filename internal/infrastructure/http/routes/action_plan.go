package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type ActionPlanRoutes struct {
	actionPlanHandler *api.ActionPlanHandler
	activityHandler   *api.ActionPlanActivityHandler
}

func NewActionPlanRoutes(actionPlanHandler *api.ActionPlanHandler, activityHandler *api.ActionPlanActivityHandler) *ActionPlanRoutes {
	return &ActionPlanRoutes{
		actionPlanHandler: actionPlanHandler,
		activityHandler:   activityHandler,
	}
}

func (r *ActionPlanRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	actionPlans := v1.Group("/action-plans")
	actionPlans.Use(middleware.PartnerMiddleware())
	{
		actionPlans.POST("", r.actionPlanHandler.Create)
		actionPlans.GET("", r.actionPlanHandler.List)

		// Rotas específicas ANTES de :id para evitar conflito
		actionPlans.GET("/department/:departmentId", r.actionPlanHandler.ListByDepartment)
		actionPlans.GET("/status", r.actionPlanHandler.ListByStatus)

		// Rotas de activities (ANTES de :id)
		actionPlans.POST("/:id/activities", r.activityHandler.Create)
		actionPlans.GET("/:id/activities", r.activityHandler.List)

		// Rotas com :id por último
		actionPlans.GET("/:id", r.actionPlanHandler.GetByID)
		actionPlans.PUT("/:id", r.actionPlanHandler.Update)
		actionPlans.DELETE("/:id", r.actionPlanHandler.Delete)
	}
}
