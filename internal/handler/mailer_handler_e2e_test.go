package handler

import (
	"bitcoin-exchange-rate/internal/repository"
	"bitcoin-exchange-rate/internal/service"
	"bitcoin-exchange-rate/pkg/mailer"
	"bitcoin-exchange-rate/pkg/parser"
	"bytes"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMailerHandler_Subscribe(t *testing.T) {
	if err := godotenv.Load("../../.env.test"); err != nil {
		t.Fatal("Failed to load .env.test file")
	}
	app := fiber.New()
	api := app.Group("/api")

	cryptoParser := parser.NewBinanceCryptoParser(os.Getenv("BASE_URL"))
	cryptoMailer := mailer.NewMailer("smtp.gmail.com", "587", os.Getenv("SENDER_EMAIL"), os.Getenv("SENDER_PASSWORD"))
	subscriberRepository := repository.NewSubscriberFileRepository(os.Getenv("TEST_FILE_PATH"))

	mailerService := service.NewMailerService(subscriberRepository, cryptoMailer)
	mailerHandler := NewMailerHandler(mailerService, cryptoParser, subscriberRepository, validator.New())

	api.Post("/subscribe", mailerHandler.Subscribe)

	tests := []struct {
		name               string
		expectedStatusCode int
		body               string
	}{
		{
			name:               "Subscribe successful",
			expectedStatusCode: fiber.StatusOK,
			body:               `{"email": "example@example.com"}`,
		},
		{
			name:               "Invalid request body",
			expectedStatusCode: fiber.StatusBadRequest,
			body:               ``,
		},
		{
			name:               "Invalid email address",
			expectedStatusCode: fiber.StatusBadRequest,
			body:               `{"email": "invalid-email"}`,
		},
		{
			name:               "Already subscribed",
			expectedStatusCode: fiber.StatusConflict,
			body:               `{"email": "example@example.com"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/subscribe", bytes.NewBufferString(test.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					t.Fatal(err)
				}
			}(resp.Body)

			require.NoError(t, err)
			require.Equal(t, test.expectedStatusCode, resp.StatusCode)
		})
	}

	if err := subscriberRepository.ClearFile(); err != nil {
		t.Fatal(err)
	}
}

func TestMailerHandler_SendEmails(t *testing.T) {
	if err := godotenv.Load("../../.env.test"); err != nil {
		t.Fatal("Failed to load .env.test file")
	}

	tests := []struct {
		name               string
		filepath           string
		expectedStatusCode int
	}{
		{
			name:               "Send emails successful",
			filepath:           os.Getenv("FULL_TEST_FILE_PATH"),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Send emails error (no subscribers)",
			filepath:           os.Getenv("EMPTY_TEST_FILE_PATH"),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			api := app.Group("/api")

			cryptoParser := parser.NewBinanceCryptoParser(os.Getenv("BASE_URL"))
			cryptoMailer := mailer.NewMailer("smtp.gmail.com", "587", os.Getenv("SENDER_EMAIL"), os.Getenv("SENDER_PASSWORD"))
			subscriberRepository := repository.NewSubscriberFileRepository(test.filepath)

			mailerService := service.NewMailerService(subscriberRepository, cryptoMailer)
			mailerHandler := NewMailerHandler(mailerService, cryptoParser, subscriberRepository, validator.New())

			api.Post("/sendEmails", mailerHandler.SendExchangeRate)

			req := httptest.NewRequest(http.MethodPost, "/api/sendEmails", nil)

			resp, err := app.Test(req, 5000)
			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					t.Fatal(err)
				}
			}(resp.Body)

			require.NoError(t, err)
			require.Equal(t, test.expectedStatusCode, resp.StatusCode)
		})
	}
}
