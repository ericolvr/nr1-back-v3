package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type PartnerHandler struct {
	partnerService *services.PartnerService
}

func NewPartnerHandler(partnerService *services.PartnerService) *PartnerHandler {
	return &PartnerHandler{
		partnerService: partnerService,
	}
}

func (h *PartnerHandler) Create(c *gin.Context) {
	var req dto.PartnerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	partner := &domain.Partner{
		Name:   req.Name,
		CNPJ:   req.CNPJ,
		Email:  req.Email,
		Mobile: req.Mobile,
		Active: true,
	}

	if err := h.partnerService.Create(c.Request.Context(), partner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar partner",
			"details": err.Error(),
		})
		return
	}

	response := h.toPartnerResponse(partner)
	c.JSON(http.StatusCreated, response)
}

func (h *PartnerHandler) List(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	partners, err := h.partnerService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar partners",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.PartnerResponse
	for _, partner := range partners {
		responses = append(responses, h.toPartnerResponse(partner))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *PartnerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	partner, err := h.partnerService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Partner não encontrado",
			"details": err.Error(),
		})
		return
	}

	response := h.toPartnerResponse(partner)
	c.JSON(http.StatusOK, response)
}

func (h *PartnerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	var updateDTO dto.PartnerRequest
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	partner := &domain.Partner{
		ID:     id,
		Name:   updateDTO.Name,
		CNPJ:   updateDTO.CNPJ,
		Email:  updateDTO.Email,
		Mobile: updateDTO.Mobile,
		Active: updateDTO.Active,
	}

	if err := h.partnerService.Update(c.Request.Context(), partner); err != nil {
		if err.Error() == "partner not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Partner não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := h.toPartnerResponse(partner)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *PartnerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	partner, err := h.partnerService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "partner not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Partner não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.partnerService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Partner deletado com sucesso",
		"data":    h.toPartnerResponse(partner),
	})
}

func (h *PartnerHandler) toPartnerResponse(partner *domain.Partner) dto.PartnerResponse {
	return dto.PartnerResponse{
		ID:        int(partner.ID),
		Name:      partner.Name,
		CNPJ:      partner.CNPJ,
		Email:     partner.Email,
		Mobile:    partner.Mobile,
		Active:    partner.Active,
		CreatedAt: partner.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
