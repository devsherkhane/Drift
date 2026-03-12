package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail is a generic helper to send notifications via SMTP
func SendEmail(toEmail string, subject string, body string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	// Standard RFC 822 email format
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, toEmail, subject, body)

	auth := smtp.PlainAuth("", from, pass, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, []byte(msg))
	return err
}

// Update your old function to use the new generic one to keep code DRY
func SendInvitationEmail(toEmail string, boardTitle string) error {
	subject := "You've been added to a board!"
	body := fmt.Sprintf("You have been added as a collaborator to the board: %s", boardTitle)
	return SendEmail(toEmail, subject, body)
}

func SendResetEmail(toEmail string, token string) error {
	subject := "Trello Clone Password Reset Request"
	body := fmt.Sprintf("Use the following token to reset your password: %s", token)
	return SendEmail(toEmail, subject, body)
}