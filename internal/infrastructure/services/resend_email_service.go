package services

import (
	"context"
	"fmt"

	"github.com/lgxju/gogretago/config"
	"github.com/lgxju/gogretago/internal/domain/services"
	"github.com/resend/resend-go/v2"
)

// ResendEmailService implements EmailService using Resend
type ResendEmailService struct {
	client    *resend.Client
	fromEmail string
}

// NewResendEmailService creates a new ResendEmailService
func NewResendEmailService() services.EmailService {
	cfg := config.Get()
	return &ResendEmailService{
		client:    resend.NewClient(cfg.ResendAPIKey),
		fromEmail: cfg.ResendFromEmail,
	}
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *ResendEmailService) SendWelcomeEmail(to, firstName string) error {
	html := fmt.Sprintf(`
		<h1>Welcome, %s!</h1>
		<p>Thank you for joining our carpooling platform.</p>
		<p>Start exploring rides and save money while reducing your carbon footprint!</p>
	`, firstName)

	return s.Send(services.SendEmailOptions{
		To:      to,
		Subject: "Welcome to Carpooling!",
		HTML:    html,
	})
}

// Send sends an email using Resend
func (s *ResendEmailService) Send(options services.SendEmailOptions) error {
	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{options.To},
		Subject: options.Subject,
		Html:    options.HTML,
	}

	_, err := s.client.Emails.SendWithContext(context.Background(), params)
	return err
}
