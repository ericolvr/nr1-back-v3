package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type AssessmentTemplateHandler struct {
	templateService *services.AssessmentTemplateService
}

func NewAssessmentTemplateHandler(templateService *services.AssessmentTemplateService) *AssessmentTemplateHandler {
	return &AssessmentTemplateHandler{
		templateService: templateService,
	}
}

func (h *AssessmentTemplateHandler) Create(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Version     *int   `json:"version"`
		Active      *bool  `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	version := 1
	if req.Version != nil {
		version = *req.Version
	}

	template := &domain.AssessmentTemplate{
		PartnerID:   partnerID,
		Name:        req.Name,
		Description: req.Description,
		Version:     version,
		Active:      active,
	}

	if err := h.templateService.Create(c.Request.Context(), template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *AssessmentTemplateHandler) List(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	templates, err := h.templateService.List(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func (h *AssessmentTemplateHandler) GetByID(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	template, err := h.templateService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *AssessmentTemplateHandler) Update(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Version     *int    `json:"version"`
		Active      *bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.templateService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Version != nil {
		template.Version = *req.Version
	}
	if req.Active != nil {
		template.Active = *req.Active
	}

	userID := c.GetInt64("user_id")
	if err := h.templateService.Update(c.Request.Context(), template, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *AssessmentTemplateHandler) Delete(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.templateService.Delete(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

func (h *AssessmentTemplateHandler) ToggleActive(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	var req dto.ToggleActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.templateService.ToggleActive(c.Request.Context(), partnerID, id, req.Active); err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Template não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	message := "Template desativado com sucesso"
	if req.Active {
		message = "Template reativado com sucesso"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h *AssessmentTemplateHandler) ListDeleted(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	templates, err := h.templateService.ListDeleted(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar templates",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, templates)
}
