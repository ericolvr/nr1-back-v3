package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type CompanyRoutes struct {
	companyHandler *api.CompanyHandler
}

func NewCompanyRoutes(companyHandler *api.CompanyHandler) *CompanyRoutes {
	return &CompanyRoutes{
		companyHandler: companyHandler,
	}
}

func (r *CompanyRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	companies := v1.Group("/companies")
	companies.Use(middleware.PartnerMiddleware())
	{
		companies.POST("", r.companyHandler.Create)
		companies.GET("", r.companyHandler.List)
		companies.GET("/deleted", r.companyHandler.ListDeleted)
		companies.GET("/:id", r.companyHandler.GetByID)
		companies.PUT("/:id", r.companyHandler.Update)
		companies.DELETE("/:id", r.companyHandler.Delete)
		companies.PATCH("/:id/toggle-active", r.companyHandler.ToggleActive)
	}
}
