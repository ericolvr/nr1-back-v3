package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type AssessmentVersionHandler struct {
	versionService *services.AssessmentVersionService
}

func NewAssessmentVersionHandler(versionService *services.AssessmentVersionService) *AssessmentVersionHandler {
	return &AssessmentVersionHandler{
		versionService: versionService,
	}
}

func (h *AssessmentVersionHandler) ListByTemplate(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	templateID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do template inválido"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	versions, err := h.versionService.ListByTemplate(c.Request.Context(), partnerID, templateID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar histórico de versões",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, versions)
}

func (h *AssessmentVersionHandler) GetByID(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	version, err := h.versionService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Versão não encontrada"})
		return
	}

	c.JSON(http.StatusOK, version)
}
