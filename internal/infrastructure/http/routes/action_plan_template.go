package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type ActionPlanTemplateRoutes struct {
	templateHandler *api.ActionPlanTemplateHandler
}

func NewActionPlanTemplateRoutes(templateHandler *api.ActionPlanTemplateHandler) *ActionPlanTemplateRoutes {
	return &ActionPlanTemplateRoutes{
		templateHandler: templateHandler,
	}
}

func (r *ActionPlanTemplateRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	templates := v1.Group("/action-plan-templates")
	templates.Use(middleware.PartnerMiddleware())
	{
		templates.POST("", r.templateHandler.Create)
		templates.GET("", r.templateHandler.List)
		templates.GET("/:id", r.templateHandler.GetByID)
		templates.PUT("/:id", r.templateHandler.Update)
		templates.DELETE("/:id", r.templateHandler.Delete)
	}
}
