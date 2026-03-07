package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type EmployeeSubmissionHandler struct {
	submissionService *services.EmployeeSubmissionService
}

func NewEmployeeSubmissionHandler(submissionService *services.EmployeeSubmissionService) *EmployeeSubmissionHandler {
	return &EmployeeSubmissionHandler{
		submissionService: submissionService,
	}
}

func (h *EmployeeSubmissionHandler) GetByToken(c *gin.Context) {
	token := c.Param("token")

	submission, err := h.submissionService.GetByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

func (h *EmployeeSubmissionHandler) Submit(c *gin.Context) {
	_ = c.Param("token") // TODO: usar token quando Submit for implementado

	var req struct {
		Answers []struct {
			QuestionID int64 `json:"question_id" binding:"required"`
			Value      int   `json:"value" binding:"required,min=1,max=5"`
		} `json:"answers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implementar método Submit no EmployeeSubmissionService
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Submit not implemented yet"})
}

func (h *EmployeeSubmissionHandler) List(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	submissions, err := h.submissionService.List(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submissions)
}
