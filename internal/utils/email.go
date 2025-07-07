package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"finsolvz-backend/internal/utils/errors"
)

type EmailService interface {
	SendForgotPasswordEmail(to, name, newPassword string) error
}

type emailService struct {
	smtpHost string
	smtpPort string
	email    string
	password string
}

func NewEmailService() EmailService {
	return &emailService{
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
		email:    os.Getenv("NODEMAILER_EMAIL"),
		password: os.Getenv("NODEMAILER_PASS"),
	}
}

func (e *emailService) SendForgotPasswordEmail(to, name, newPassword string) error {
	if e.email == "" || e.password == "" {
		return errors.New("EMAIL_CONFIG_MISSING", "Email configuration not found", 500, nil, nil)
	}

	// Email template
	emailTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset - Finsolvz</title>
</head>
<body style="font-family: sans-serif; line-height: 1.6; margin: 0; padding: 20px;">
    <div style="max-width: 600px; margin: 0 auto;">
        <h2>Password Reset - Finsolvz</h2>
        <p>Dear <strong>{{.Name}}</strong>,</p>
        <p>We have received a request to reset your password for your <strong>Finsolvz</strong> account.</p>
        <p>Here is your new password:</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <p style="font-size: 18px; font-weight: bold; margin: 0; font-family: monospace;">{{.NewPassword}}</p>
        </div>
        <p>Please use this password to log in to your account. For security reasons, we strongly recommend changing your password after logging in.</p>
        <p>If you did not request this change, please contact our support team immediately.</p>
        <p style="margin-top: 30px;">Best regards,<br/>Finsolvz Team</p>
    </div>
</body>
</html>`

	// Parse template
	tmpl, err := template.New("forgotPassword").Parse(emailTemplate)
	if err != nil {
		return errors.New("EMAIL_TEMPLATE_ERROR", "Failed to parse email template", 500, err, nil)
	}

	// Execute template
	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		Name        string
		NewPassword string
	}{
		Name:        name,
		NewPassword: newPassword,
	})
	if err != nil {
		return errors.New("EMAIL_TEMPLATE_ERROR", "Failed to execute email template", 500, err, nil)
	}

	// Compose email
	subject := "Your New Finsolvz Account Password"
	message := fmt.Sprintf("From: Finsolvz <%s>\r\n", e.email)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body.String()

	// Send email
	auth := smtp.PlainAuth("", e.email, e.password, e.smtpHost)
	err = smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.email, []string{to}, []byte(message))
	if err != nil {
		return errors.New("EMAIL_SEND_ERROR", "Failed to send email", 500, err, nil)
	}

	return nil
}
