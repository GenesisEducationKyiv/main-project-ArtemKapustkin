package model

type EmailMessage string

func NewEmailMessage(message string) EmailMessage {
	return EmailMessage(message)
}
