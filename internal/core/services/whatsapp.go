package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type WhatsAppService struct {
	accountSID string
	authToken  string
	fromNumber string
}

func NewWhatsAppService() *WhatsAppService {
	return &WhatsAppService{
		accountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		authToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		fromNumber: os.Getenv("TWILIO_WHATSAPP_FROM"),
	}
}

// SendInvitation envia convite via WhatsApp
func (s *WhatsAppService) SendInvitation(toPhone, templateName, invitationToken, surveyURL string) error {
	// Formatar número de telefone para formato internacional
	toPhone = s.formatPhoneNumber(toPhone)

	message := fmt.Sprintf(
		"*Convite para Avaliação NR-1*\n\n"+
			"Você foi convidado(a) para participar da avaliação:\n\n"+
			"📋 *%s*\n\n"+
			"Sua participação é muito importante! Clique no link abaixo para começar:\n\n"+
			"%s\n\n"+

			templateName,
		surveyURL,
		invitationToken,
	)

	return s.sendMessage(toPhone, message)
}

// formatPhoneNumber formata número de telefone para formato internacional do WhatsApp
func (s *WhatsAppService) formatPhoneNumber(phone string) string {
	// Remover espaços e caracteres especiais
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Se já tem whatsapp: prefix, remover para processar
	phone = strings.TrimPrefix(phone, "whatsapp:")

	// Se não tem +, adicionar +55 (Brasil) se começar com dígito
	if !strings.HasPrefix(phone, "+") && len(phone) > 0 && phone[0] >= '0' && phone[0] <= '9' {
		phone = "+55" + phone
	}

	// Adicionar whatsapp: prefix
	return "whatsapp:" + phone
}

// SendNotification envia notificação genérica via WhatsApp
func (s *WhatsAppService) SendNotification(toPhone, title, message string) error {
	// Formatar número de telefone para WhatsApp
	toPhone = s.formatPhoneNumber(toPhone)

	formattedMessage := fmt.Sprintf(
		"🔔 *%s*\n\n%s",
		title,
		message,
	)

	return s.sendMessage(toPhone, formattedMessage)
}

// sendMessage envia mensagem via Twilio WhatsApp API
func (s *WhatsAppService) sendMessage(to, body string) error {
	fmt.Printf("🔧 [WhatsApp] Validando configurações...\n")

	if s.accountSID == "" || s.authToken == "" || s.fromNumber == "" {
		fmt.Printf("❌ [WhatsApp] Configurações incompletas:\n")
		fmt.Printf("   - Account SID: %s\n", func() string {
			if s.accountSID == "" {
				return "VAZIO"
			}
			return "***"
		}())
		fmt.Printf("   - Auth Token: %s\n", func() string {
			if s.authToken == "" {
				return "VAZIO"
			}
			return "***"
		}())
		fmt.Printf("   - From Number: %s\n", s.fromNumber)
		return fmt.Errorf("configurações do Twilio WhatsApp incompletas. Verifique as variáveis de ambiente")
	}

	fmt.Printf("✅ [WhatsApp] Configurações OK\n")
	fmt.Printf("   - Account SID: %s\n", s.accountSID)
	fmt.Printf("   - Auth Token: %s...\n", s.authToken[:8])
	fmt.Printf("   - From: %s\n", s.fromNumber)

	// Twilio API endpoint
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.accountSID)

	fmt.Printf("📨 [WhatsApp] Montando mensagem:\n")
	fmt.Printf("   - De: %s\n", s.fromNumber)
	fmt.Printf("   - Para: %s\n", to)
	fmt.Printf("   - Mensagem: %d caracteres\n", len(body))

	// Preparar dados do formulário
	data := url.Values{}
	data.Set("From", s.fromNumber)
	data.Set("To", to)
	data.Set("Body", body)

	// Criar request
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("❌ [WhatsApp] Erro ao criar request: %v\n", err)
		return fmt.Errorf("erro ao criar request: %v", err)
	}

	// Adicionar headers
	req.SetBasicAuth(s.accountSID, s.authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Printf("🚀 [WhatsApp] Enviando mensagem via Twilio...\n")

	// Enviar request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ [WhatsApp] Falha ao enviar: %v\n", err)
		return fmt.Errorf("erro ao enviar mensagem WhatsApp: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Ler corpo da resposta para debug
		bodyBytes := make([]byte, 500)
		n, _ := resp.Body.Read(bodyBytes)
		responseBody := string(bodyBytes[:n])

		fmt.Printf("❌ [WhatsApp] Resposta com erro: %d\n", resp.StatusCode)
		fmt.Printf("   - Resposta Twilio: %s\n", responseBody)

		if resp.StatusCode == 401 {
			fmt.Printf("   - ⚠️  ERRO DE AUTENTICAÇÃO: Verifique TWILIO_ACCOUNT_SID e TWILIO_AUTH_TOKEN\n")
		}

		return fmt.Errorf("erro ao enviar mensagem WhatsApp: status %d", resp.StatusCode)
	}

	fmt.Printf("✅ [WhatsApp] Mensagem enviada com sucesso!\n")
	return nil
}

// IsConfigured verifica se o serviço está configurado
func (s *WhatsAppService) IsConfigured() bool {
	return s.accountSID != "" &&
		s.authToken != "" &&
		s.fromNumber != ""
}

// GetConfig retorna configurações (para debug)
func (s *WhatsAppService) GetConfig() map[string]string {
	return map[string]string{
		"account_sid_configured": fmt.Sprintf("%v", s.accountSID != ""),
		"auth_token_configured":  fmt.Sprintf("%v", s.authToken != ""),
		"from_number":            s.fromNumber,
		"configured":             fmt.Sprintf("%v", s.IsConfigured()),
	}
}
