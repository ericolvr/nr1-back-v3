package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	companyService *services.CompanyService
}

func NewCompanyHandler(companyService *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) Create(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	var req dto.CompanyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	company := &domain.Company{
		PartnerID: partnerID,
		Name:      req.Name,
		CNPJ:      req.CNPJ,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Active:    true,
	}

	if err := h.companyService.Create(c.Request.Context(), company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar company",
			"details": err.Error(),
		})
		return
	}

	response := h.toCompanyResponse(company)
	c.JSON(http.StatusCreated, response)
}

func (h *CompanyHandler) List(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	companies, err := h.companyService.List(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar companies",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.CompanyResponse
	for _, company := range companies {
		responses = append(responses, h.toCompanyResponse(company))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *CompanyHandler) GetByID(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	company, err := h.companyService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Company não encontrada",
			"details": err.Error(),
		})
		return
	}

	response := h.toCompanyResponse(company)
	c.JSON(http.StatusOK, response)
}

func (h *CompanyHandler) Update(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	var updateDTO dto.CompanyRequest
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	company := &domain.Company{
		ID:        id,
		PartnerID: partnerID,
		Name:      updateDTO.Name,
		CNPJ:      updateDTO.CNPJ,
		Email:     updateDTO.Email,
		Mobile:    updateDTO.Mobile,
		Active:    updateDTO.Active,
	}

	if err := h.companyService.Update(c.Request.Context(), company); err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company não encontrada",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := h.toCompanyResponse(company)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *CompanyHandler) Delete(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	company, err := h.companyService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company não encontrada",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.companyService.Delete(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company deletada com sucesso",
		"data":    h.toCompanyResponse(company),
	})
}

func (h *CompanyHandler) toCompanyResponse(company *domain.Company) dto.CompanyResponse {
	return dto.CompanyResponse{
		ID:        int(company.ID),
		Name:      company.Name,
		CNPJ:      company.CNPJ,
		Email:     company.Email,
		Mobile:    company.Mobile,
		Active:    company.Active,
		CreatedAt: company.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
