package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     os.Getenv("SMTP_PORT"),
		smtpUser:     os.Getenv("SMTP_USER"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
		fromName:     os.Getenv("SMTP_FROM_NAME"),
	}
}

func (s *EmailService) SendInvitation(toEmail, templateName, invitationToken, surveyURL string) error {
	subject := fmt.Sprintf("Convite: %s", templateName)

	// Template HTML do email - Moderno com gradiente roxo
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;
            line-height: 1.6; 
            color: #1F2937;
            margin: 0;
            padding: 0;
            background-color: #F9FAFB;
        }
        .container { 
            max-width: 600px; 
            margin: 40px auto; 
            background-color: #ffffff;
            border-radius: 16px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header { 
            background: linear-gradient(135deg, #4f38a8 0%, #8c2e9b 50%, #d53f8c 100%);
            color: white; 
            padding: 40px 30px; 
            text-align: center; 
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
            font-weight: bold;
        }
        .content { 
            padding: 40px 30px;
            background-color: #ffffff;
        }
        .content p {
            margin: 0 0 16px 0;
            color: #374151;
            font-size: 15px;
        }
        .template-name {
            background: linear-gradient(135deg, #F3E8FF 0%, #FCE7F3 100%);
            padding: 20px;
            border-radius: 12px;
            margin: 24px 0;
            border-left: 4px solid #8c2e9b;
        }
        .template-name h2 {
            margin: 0;
            color: #1F2937;
            font-size: 20px;
            font-weight: bold;
        }
        .button-container {
            text-align: center;
            margin: 32px 0;
        }
        .button { 
            display: inline-block; 
            padding: 16px 40px; 
            background: linear-gradient(135deg, #4f38a8 0%, #8c2e9b 50%, #d53f8c 100%);
            color: white !important; 
            text-decoration: none; 
            border-radius: 12px;
            font-weight: bold;
            font-size: 16px;
            box-shadow: 0 4px 12px rgba(140, 46, 155, 0.3);
        }
        .button:hover {
            opacity: 0.9;
        }
        .link-box {
            background-color: #F9FAFB;
            padding: 16px;
            border-radius: 8px;
            margin: 24px 0;
            border: 1px solid #E5E7EB;
        }
        .link-box p {
            margin: 0 0 8px 0;
            font-size: 13px;
            color: #6B7280;
        }
        .link-box a {
            color: #8c2e9b;
            word-break: break-all;
            font-size: 13px;
        }
        .token-box {
            background-color: #FEF3C7;
            border: 1px solid #FCD34D;
            padding: 16px;
            border-radius: 8px;
            margin: 24px 0;
        }
        .token-box p {
            margin: 0;
            font-size: 13px;
            color: #92400E;
        }
        .token-box strong {
            color: #78350F;
            font-size: 14px;
        }
        .footer { 
            text-align: center; 
            padding: 24px 30px;
            background-color: #F9FAFB;
            border-top: 1px solid #E5E7EB;
        }
        .footer p {
            margin: 0;
            font-size: 12px; 
            color: #6B7280;
        }
        @media only screen and (max-width: 600px) {
            .container {
                margin: 0;
                border-radius: 0;
            }
            .header {
                padding: 30px 20px;
            }
            .header h1 {
                font-size: 24px;
            }
            .content {
                padding: 30px 20px;
            }
            .button {
                padding: 14px 32px;
                font-size: 15px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Convite para Avaliação NR-1</h1>
        </div>
        
        <div class="content">
            <p><strong>Olá,</strong></p>
            
            <p>Você foi convidado(a) para participar da avaliação:</p>
            
            <div class="template-name">
                <h2>{{.TemplateName}}</h2>
            </div>
            
            <p>Sua participação é muito importante para nós! Suas respostas nos ajudarão a melhorar continuamente nosso ambiente de trabalho e garantir a conformidade com a NR-1.</p>
            
            <p><strong>Clique no botão abaixo para começar:</strong></p>
            
            <div class="button-container">
                <a href="{{.SurveyURL}}" class="button">Responder Avaliação</a>
            </div>
            
            <div class="link-box">
                <p><strong>Ou copie e cole este link no seu navegador:</strong></p>
                <a href="{{.SurveyURL}}">{{.SurveyURL}}</a>
            </div>
            
            <div class="token-box">
                <p><strong>Código de acesso:</strong> {{.InvitationToken}}</p>
            </div>
        </div>
        
        <div class="footer">
            <p>Este é um email automático. Por favor, não responda.</p>
            <p style="margin-top: 8px;">© 2024 - Todos os direitos reservados</p>
        </div>
    </div>
</body>
</html>
`

	tmpl, err := template.New("invitation").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("erro ao criar template: %v", err)
	}

	var body bytes.Buffer
	data := struct {
		TemplateName    string
		SurveyURL       string
		InvitationToken string
	}{
		TemplateName:    templateName,
		SurveyURL:       surveyURL,
		InvitationToken: invitationToken,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("erro ao executar template: %v", err)
	}

	return s.sendEmail(toEmail, subject, body.String())
}

func (s *EmailService) SendTestEmail(toEmail string) error {
	subject := "Teste de Configuração SMTP"
	body := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h2>✅ Configuração SMTP Funcionando!</h2>
    <p>Este é um email de teste para verificar se o serviço de email está configurado corretamente.</p>
    <p>Se você recebeu este email, significa que o SMTP está funcionando perfeitamente.</p>
</body>
</html>
`
	return s.sendEmail(toEmail, subject, body)
}

func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	fmt.Printf("🔧 [SMTP] Validando configurações...\n")

	if s.smtpHost == "" || s.smtpPort == "" || s.smtpUser == "" || s.smtpPassword == "" {
		fmt.Printf("❌ [SMTP] Configurações incompletas:\n")
		fmt.Printf("   - Host: %s\n", s.smtpHost)
		fmt.Printf("   - Port: %s\n", s.smtpPort)
		fmt.Printf("   - User: %s\n", s.smtpUser)
		fmt.Printf("   - Password: %s\n", func() string {
			if s.smtpPassword == "" {
				return "VAZIO"
			}
			return "***"
		}())
		return fmt.Errorf("configurações SMTP incompletas. Verifique as variáveis de ambiente")
	}

	fmt.Printf("✅ [SMTP] Configurações OK\n")
	fmt.Printf("   - Servidor: %s:%s\n", s.smtpHost, s.smtpPort)
	fmt.Printf("   - Usuário: %s\n", s.smtpUser)

	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)

	from := s.fromEmail
	if from == "" {
		from = s.smtpUser
	}

	fromName := s.fromName
	if fromName == "" {
		fromName = "NR-1 Sistema"
	}

	fmt.Printf("📨 [SMTP] Montando mensagem:\n")
	fmt.Printf("   - De: %s <%s>\n", fromName, from)
	fmt.Printf("   - Para: %s\n", to)
	fmt.Printf("   - Assunto: %s\n", subject)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	fmt.Printf("🚀 [SMTP] Conectando em %s...\n", addr)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		fmt.Printf("❌ [SMTP] Falha ao enviar: %v\n", err)
		return fmt.Errorf("erro ao enviar email: %v", err)
	}

	fmt.Printf("✅ [SMTP] Email enviado com sucesso!\n")
	return nil
}

func (s *EmailService) IsConfigured() bool {
	return s.smtpHost != "" &&
		s.smtpPort != "" &&
		s.smtpUser != "" &&
		s.smtpPassword != ""
}

func (s *EmailService) GetConfig() map[string]string {
	return map[string]string{
		"smtp_host":       s.smtpHost,
		"smtp_port":       s.smtpPort,
		"smtp_user":       s.smtpUser,
		"smtp_configured": fmt.Sprintf("%v", s.IsConfigured()),
		"from_email":      s.fromEmail,
		"from_name":       s.fromName,
	}
}
