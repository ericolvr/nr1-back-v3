package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type ActionPlanRoutes struct {
	actionPlanHandler *api.ActionPlanHandler
}

func NewActionPlanRoutes(actionPlanHandler *api.ActionPlanHandler) *ActionPlanRoutes {
	return &ActionPlanRoutes{
		actionPlanHandler: actionPlanHandler,
	}
}

func (r *ActionPlanRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	actionPlans := v1.Group("/action-plans")
	actionPlans.Use(middleware.PartnerMiddleware())
	{
		actionPlans.POST("", r.actionPlanHandler.Create)
		actionPlans.GET("", r.actionPlanHandler.List)
		actionPlans.GET("/:id", r.actionPlanHandler.GetByID)
		actionPlans.PUT("/:id", r.actionPlanHandler.Update)
		actionPlans.DELETE("/:id", r.actionPlanHandler.Delete)
		actionPlans.GET("/department/:departmentId", r.actionPlanHandler.ListByDepartment)
		actionPlans.GET("/status", r.actionPlanHandler.ListByStatus)
	}
}
