package utils

import (
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/smtp"
	"os"
)

func SendEmailWithSendGrid(to, cc, bcc, textContent, htmlBody string) error {
	from := mail.NewEmail("File Sender", os.Getenv("SENDGRID_FROM"))
	toEmail := mail.NewEmail("Recipient", to)
	var subject = "Your Secure File Link Sent By" + "...."
	message := mail.NewSingleEmail(from, subject, toEmail, textContent, htmlBody)

	// Add CC and BCC if provided
	if cc != "" {
		message.Personalizations[0].AddCCs(mail.NewEmail("CC", cc))
	}
	if bcc != "" {
		message.Personalizations[0].AddBCCs(mail.NewEmail("BCC", bcc))
	}

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid error: %v - %v", response.StatusCode, response.Body)
	}

	return nil
}

func SendEmail(toEmail, bccEmail, ccEmail, body string) error {
	e := email.NewEmail()
	e.From = "file <" + os.Getenv("SMTP_USERNAME") + ">"

	if toEmail != "" {
		e.To = []string{toEmail}
	} else {
		return errors.New("email address is required")
	}

	// Set Bcc only if non-empty
	if bccEmail != "" {
		e.Bcc = []string{bccEmail}
	}

	// Set Cc only if non-empty
	if ccEmail != "" {
		e.Cc = []string{ccEmail}
	}

	e.Subject = "Your Secure File Link Sent By" + "...."
	e.Text = []byte(body)
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")

	return e.Send(os.Getenv("SMTP_ADDRESS"), smtp.PlainAuth(
		"",
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"),
	))
}
