package notification

import (
	"fmt"
	"managify/internal/config"
	"net/smtp"
)

// Provider defines the interface for sending notifications.
type Provider interface {
	SendVerificationEmail(to, token string) error
}

type smtpProvider struct {
	cfg *config.Config
}

// NewSMTPProvider creates a new SMTP-based notification provider.
func NewSMTPProvider(cfg *config.Config) Provider {
	return &smtpProvider{cfg: cfg}
}

func (p *smtpProvider) SendVerificationEmail(to, token string) error {
	subject := "Subject: Email Verification\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	
	// Use FrontendURL from config instead of hardcoded localhost
	verificationURL := fmt.Sprintf("%s/verify?token=%s", p.cfg.FrontendURL, token)
	
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Verify Your Email</h2>
			<p>Click the button below to verify your account:</p>
			<a href='%s' 
			style='display:inline-block;padding:10px 20px;background-color:#4CAF50;color:white;text-decoration:none;border-radius:5px;'>Verify Email</a>
			<p>If you did not create an account, you can ignore this email.</p>
		</body>
		</html>`, verificationURL)

	msg := subject + mime + body
	addr := fmt.Sprintf("%s:%s", p.cfg.SMTPHost, p.cfg.SMTPPort)
	auth := smtp.PlainAuth("", p.cfg.SMTPUser, p.cfg.SMTPPass, p.cfg.SMTPHost)

	return smtp.SendMail(addr, auth, p.cfg.SMTPFrom, []string{to}, []byte(msg))
}
