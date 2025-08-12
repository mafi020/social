package utils

import (
	"fmt"
	"log"

	"github.com/mafi020/social/internal/env"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(fromName, fromEmail, toName, toEmail, subject, plainTextContent, htmlContent string) error {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(env.GetEnvOrPanic("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent: Status %d \n", response.StatusCode)
	return nil
}
