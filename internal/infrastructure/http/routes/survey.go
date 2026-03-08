package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type SurveyRoutes struct {
	handler *api.SurveyHandler
}

func NewSurveyRoutes(handler *api.SurveyHandler) *SurveyRoutes {
	return &SurveyRoutes{handler: handler}
}

func (r *SurveyRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	// Rota pública - SEM middleware de autenticação
	// Permite acesso anônimo via token do convite
	v1.GET("/survey", r.handler.GetSurvey)
}
