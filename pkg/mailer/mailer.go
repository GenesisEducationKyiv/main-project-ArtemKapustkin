package mailer

import (
	"net/smtp"
)

type Mailer struct {
	server   string
	port     string
	email    string
	password string
	auth     smtp.Auth
}

func NewMailer(server, port, email, password string) *Mailer {
	return &Mailer{
		server:   server,
		port:     port,
		email:    email,
		password: password,
		auth:     smtp.PlainAuth("", email, password, server),
	}
}

func (m *Mailer) SendEmail(email, message string) error {
	err := smtp.SendMail(m.server+":"+m.port, m.auth, m.email, []string{email}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
