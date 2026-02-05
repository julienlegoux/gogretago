package services

// SendEmailOptions contains options for sending an email
type SendEmailOptions struct {
	To      string
	Subject string
	HTML    string
}

// EmailService defines the interface for email operations
type EmailService interface {
	SendWelcomeEmail(to string, firstName string) error
	Send(options SendEmailOptions) error
}
