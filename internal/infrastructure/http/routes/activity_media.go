package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type ActivityMediaRoutes struct {
	mediaHandler *api.ActivityMediaHandler
}

func NewActivityMediaRoutes(mediaHandler *api.ActivityMediaHandler) *ActivityMediaRoutes {
	return &ActivityMediaRoutes{
		mediaHandler: mediaHandler,
	}
}

func (r *ActivityMediaRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	// Rotas de upload de mídia para atividades
	activities := v1.Group("/activities")
	activities.Use(middleware.PartnerMiddleware())
	{
		// Upload de arquivo
		activities.POST("/:activityId/media", r.mediaHandler.Upload)
		// Listar mídias de uma atividade
		activities.GET("/:activityId/media", r.mediaHandler.List)
	}

	// Rotas diretas para mídias
	media := v1.Group("/activity-media")
	media.Use(middleware.PartnerMiddleware())
	{
		// Deletar mídia
		media.DELETE("/:id", r.mediaHandler.Delete)
	}
}
