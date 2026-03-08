package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type EmailRoutes struct {
	handler *api.EmailHandler
}

func NewEmailRoutes(handler *api.EmailHandler) *EmailRoutes {
	return &EmailRoutes{handler: handler}
}

func (r *EmailRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	email := v1.Group("/email")
	email.Use(middleware.PartnerMiddleware())
	{
		email.GET("/config", r.handler.CheckConfig)
		email.POST("/test", r.handler.TestEmail)
		email.POST("/send-invitation", r.handler.SendInvitation)
	}
}
