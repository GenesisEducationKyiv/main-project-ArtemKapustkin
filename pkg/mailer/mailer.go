package mailer

import (
	"net/smtp"
)

type Mailer struct {
	server   string
	port     string
	email    string
	password string
}

func NewMailer(server, port, email, password string) *Mailer {
	return &Mailer{
		server:   server,
		port:     port,
		email:    email,
		password: password,
	}
}

func (m *Mailer) SendEmail(email, message string) error {
	auth := smtp.PlainAuth("", m.email, m.password, m.server)

	err := smtp.SendMail(m.server+":"+m.port, auth, m.email, []string{email}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
