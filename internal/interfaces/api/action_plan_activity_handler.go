package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type ActionPlanActivityHandler struct {
	activityRepo domain.ActionPlanActivityRepository
}

func NewActionPlanActivityHandler(activityRepo domain.ActionPlanActivityRepository) *ActionPlanActivityHandler {
	return &ActionPlanActivityHandler{
		activityRepo: activityRepo,
	}
}

type CreateActivityRequest struct {
	Type          string  `json:"type" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
	DueDate       *string `json:"due_date"`
	CreatedBy     *int64  `json:"created_by"`
	CreatedByName string  `json:"created_by_name"`
}

func (h *ActionPlanActivityHandler) Create(c *gin.Context) {
	actionPlanID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action plan ID"})
		return
	}

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	activity := &domain.ActionPlanActivity{
		ActionPlanID:  actionPlanID,
		Type:          req.Type,
		Title:         req.Title,
		Description:   req.Description,
		Status:        req.Status,
		CreatedBy:     req.CreatedBy,
		CreatedByName: req.CreatedByName,
	}

	if err := activity.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := h.activityRepo.Create(c.Request.Context(), activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func (h *ActionPlanActivityHandler) List(c *gin.Context) {
	actionPlanID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action plan ID"})
		return
	}

	activities, err := h.activityRepo.ListByActionPlan(c.Request.Context(), actionPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list activities", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (h *ActionPlanActivityHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := h.activityRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func (h *ActionPlanActivityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Buscar activity existente
	activity, err := h.activityRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	// Atualizar campos
	activity.Title = req.Title
	activity.Description = req.Description
	activity.Status = req.Status

	if err := activity.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := h.activityRepo.Update(c.Request.Context(), activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func (h *ActionPlanActivityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	if err := h.activityRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

func (h *ActionPlanActivityHandler) MarkAsCompleted(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := h.activityRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	activity.MarkAsCompleted()

	if err := h.activityRepo.Update(c.Request.Context(), activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activity)
}
