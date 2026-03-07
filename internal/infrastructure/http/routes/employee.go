package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type EmployeeRoutes struct {
	employeeHandler *api.EmployeeHandler
}

func NewEmployeeRoutes(employeeHandler *api.EmployeeHandler) *EmployeeRoutes {
	return &EmployeeRoutes{
		employeeHandler: employeeHandler,
	}
}

func (er *EmployeeRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	employees := v1.Group("/employees")
	employees.Use(middleware.PartnerMiddleware())
	// employees.Use(middleware.JWTMiddleware())
	{
		employees.POST("", er.employeeHandler.Create)
		employees.GET("", er.employeeHandler.List)
		employees.GET("/deleted", er.employeeHandler.ListDeleted)
		employees.GET("/:id", er.employeeHandler.GetByID)
		employees.PATCH("/:id", er.employeeHandler.Update)
		employees.DELETE("/:id", er.employeeHandler.Delete)
		employees.PATCH("/:id/toggle-active", er.employeeHandler.ToggleActive)
	}
}
