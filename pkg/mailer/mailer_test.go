package mailer

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestSendingEmail(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		t.Errorf("error loading .env.test file")
	}

	mailer := NewMailer("smtp.gmail.com", "587", os.Getenv("SENDER_EMAIL"), os.Getenv("SENDER_PASSWORD"))

	err = mailer.SendEmail("a.kapustkin.2003@gmail.com", "test_message")
	if err != nil {
		t.Errorf("failure occurs while sending email: %v", err)
	}
}
