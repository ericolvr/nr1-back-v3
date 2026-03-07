package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type PartnerRoutes struct {
	partnerHandler *api.PartnerHandler
}

func NewPartnerRoutes(partnerHandler *api.PartnerHandler) *PartnerRoutes {
	return &PartnerRoutes{
		partnerHandler: partnerHandler,
	}
}

func (r *PartnerRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	partners := v1.Group("/partners")
	{
		partners.POST("", r.partnerHandler.Create)
		partners.GET("", r.partnerHandler.List)
		partners.GET("/:id", r.partnerHandler.GetByID)
		partners.PUT("/:id", r.partnerHandler.Update)
		partners.DELETE("/:id", r.partnerHandler.Delete)
	}
}
