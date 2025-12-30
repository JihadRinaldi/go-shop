package notifications

import (
	"fmt"
	"net"
	"net/smtp"
	"strconv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailConfig struct {
	To      string
	Subject string
	Body    string
}

type EmailNotifier struct {
	config *SMTPConfig
}

func NewEmailNotifier(config *SMTPConfig) *EmailNotifier {
	return &EmailNotifier{config: config}
}

func (n *EmailNotifier) SendEmail(emailConfig *EmailConfig) error {
	addr := net.JoinHostPort(n.config.Host, strconv.Itoa(n.config.Port))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, n.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if n.config.Username != "" && n.config.Password != "" {
		auth := smtp.PlainAuth("", n.config.Username, n.config.Password, n.config.Host)
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(n.config.From); err != nil {
		return err
	}

	if err := client.Rcpt(emailConfig.To); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		n.config.From, emailConfig.To, emailConfig.Subject, emailConfig.Body)

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	return w.Close()
}

func (e *EmailNotifier) SendLoginNotification(userEmail, userName string) error {
	email := &EmailConfig{
		To:      userEmail,
		Subject: "Login Notification",
		Body: fmt.Sprintf(`Hello %s,

You have successfully logged into your account.

If this wasn't you, please contact support immediately.

Best regards,
The Shop Team`, userName),
	}

	return e.SendEmail(email)
}
