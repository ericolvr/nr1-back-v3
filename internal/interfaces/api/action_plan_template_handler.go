package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type ActionPlanTemplateHandler struct {
	templateRepo domain.ActionPlanTemplateRepository
}

func NewActionPlanTemplateHandler(templateRepo domain.ActionPlanTemplateRepository) *ActionPlanTemplateHandler {
	return &ActionPlanTemplateHandler{
		templateRepo: templateRepo,
	}
}

type CreateTemplateRequest struct {
	Category            string `json:"category" binding:"required"`
	MinRiskLevel        string `json:"min_risk_level" binding:"required"`
	TitleTemplate       string `json:"title_template" binding:"required"`
	DescriptionTemplate string `json:"description_template" binding:"required"`
	Priority            string `json:"priority" binding:"required"`
	DefaultDueDays      int    `json:"default_due_days" binding:"required"`
	AutoCreate          bool   `json:"auto_create"`
	Active              bool   `json:"active"`
}

func (h *ActionPlanTemplateHandler) Create(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	template := &domain.ActionPlanTemplate{
		PartnerID:           partnerID,
		Category:            req.Category,
		MinRiskLevel:        req.MinRiskLevel,
		TitleTemplate:       req.TitleTemplate,
		DescriptionTemplate: req.DescriptionTemplate,
		Priority:            req.Priority,
		DefaultDueDays:      req.DefaultDueDays,
		AutoCreate:          req.AutoCreate,
		Active:              req.Active,
	}

	if err := template.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := h.templateRepo.Create(c.Request.Context(), template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *ActionPlanTemplateHandler) List(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	// Check if filtering by active status
	if c.Query("active") == "true" {
		templates, err := h.templateRepo.ListActive(c.Request.Context(), partnerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list templates", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, templates)
		return
	}

	templates, err := h.templateRepo.List(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list templates", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func (h *ActionPlanTemplateHandler) GetByID(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.templateRepo.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *ActionPlanTemplateHandler) Update(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	template := &domain.ActionPlanTemplate{
		ID:                  id,
		PartnerID:           partnerID,
		Category:            req.Category,
		MinRiskLevel:        req.MinRiskLevel,
		TitleTemplate:       req.TitleTemplate,
		DescriptionTemplate: req.DescriptionTemplate,
		Priority:            req.Priority,
		DefaultDueDays:      req.DefaultDueDays,
		AutoCreate:          req.AutoCreate,
		Active:              req.Active,
	}

	if err := template.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := h.templateRepo.Update(c.Request.Context(), template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *ActionPlanTemplateHandler) Delete(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	if err := h.templateRepo.Delete(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}
