package main

import (
	"bitcoin-exchange-rate/internal/handler"
	"bitcoin-exchange-rate/internal/repository"
	"bitcoin-exchange-rate/internal/service"
	"bitcoin-exchange-rate/pkg/mailer"
	"bitcoin-exchange-rate/pkg/parser"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cryptoParser := parser.NewBinanceCryptoParser()
	cryptoMailer := mailer.NewMailer("smtp.gmail.com", "587", os.Getenv("SENDER_EMAIL"), os.Getenv("SENDER_PASSWORD"))

	subscriberRepository := repository.NewSubscriberFileRepository(os.Getenv("EMAILS_FILEPATH"))

	mailerService := service.NewMailerService(subscriberRepository, cryptoMailer)

	rateHandler := handler.NewRateHandler(cryptoParser)
	mailerHandler := handler.NewMailerHandler(mailerService, cryptoParser, subscriberRepository, validator.New())

	app := fiber.New()
	api := app.Group("/api")

	api.Get("/rate", rateHandler.GetExchangeRate)
	api.Post("/subscribe", mailerHandler.Subscribe)
	api.Post("/sendEmails", mailerHandler.SendExchangeRate)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
