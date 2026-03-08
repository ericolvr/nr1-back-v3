package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type InvitationRoutes struct {
	handler *api.InvitationHandler
}

func NewInvitationRoutes(handler *api.InvitationHandler) *InvitationRoutes {
	return &InvitationRoutes{handler: handler}
}

func (r *InvitationRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	invitations := v1.Group("/invitations")
	invitations.Use(middleware.PartnerMiddleware())
	{
		invitations.GET("", r.handler.List)
		invitations.GET("/summary", r.handler.GetSummary)
		invitations.POST("/send-all", r.handler.SendAllInvitations)
		invitations.GET("/:id", r.handler.GetByID)
		invitations.PUT("/:id/sent", r.handler.MarkAsSent)
		invitations.PUT("/:id/failed", r.handler.MarkAsFailed)
		invitations.DELETE("/:id", r.handler.Delete)
	}
}
