package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

func (h *EmailHandler) TestEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email inválido",
		})
		return
	}

	if !h.emailService.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":  "Serviço de email não configurado",
			"config": h.emailService.GetConfig(),
		})
		return
	}

	err := h.emailService.SendTestEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao enviar email de teste",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email de teste enviado com sucesso",
		"to":      req.Email,
	})
}

func (h *EmailHandler) CheckConfig(c *gin.Context) {
	config := h.emailService.GetConfig()
	isConfigured := h.emailService.IsConfigured()

	status := http.StatusOK
	if !isConfigured {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"configured": isConfigured,
		"config":     config,
	})
}

func (h *EmailHandler) SendInvitation(c *gin.Context) {
	var req struct {
		Email           string `json:"email" binding:"required,email"`
		TemplateName    string `json:"template_name" binding:"required"`
		InvitationToken string `json:"invitation_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ [EMAIL] Erro ao fazer bind do JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos",
		})
		return
	}

	fmt.Printf("📧 [EMAIL] Iniciando envio de convite:\n")
	fmt.Printf("   - Email: %s\n", req.Email)
	fmt.Printf("   - Template: %s\n", req.TemplateName)
	fmt.Printf("   - Token: %s\n", req.InvitationToken)

	if !h.emailService.IsConfigured() {
		fmt.Printf("❌ [EMAIL] Serviço de email não configurado\n")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Serviço de email não configurado",
		})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	surveyURL := fmt.Sprintf("%s/survey?token=%s", frontendURL, req.InvitationToken)
	fmt.Printf("   - URL da avaliação: %s\n", surveyURL)

	fmt.Printf("📤 [EMAIL] Enviando email...\n")
	err := h.emailService.SendInvitation(
		req.Email,
		req.TemplateName,
		req.InvitationToken,
		surveyURL,
	)

	if err != nil {
		fmt.Printf("❌ [EMAIL] Erro ao enviar: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao enviar convite",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("✅ [EMAIL] Email enviado com sucesso para %s\n", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Convite enviado com sucesso",
		"to":      req.Email,
	})
}
