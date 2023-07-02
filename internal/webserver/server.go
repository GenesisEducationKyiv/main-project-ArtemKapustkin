package webserver

import (
	"bitcoin-exchange-rate/internal/handler"
	"bitcoin-exchange-rate/internal/model"
	"bitcoin-exchange-rate/internal/repository"
	"bitcoin-exchange-rate/internal/service"
	"bitcoin-exchange-rate/pkg/mailer"
	"bitcoin-exchange-rate/pkg/parser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Config struct {
	CryptoParserBaseURL                string
	CryptoMailerSenderEmail            string
	CryptoMailerSenderPassword         string
	SubscriberRepositoryEmailsFilePath string
	BaseCurrencyStr                    string
	QuoteCurrencyStr                   string
}

type App struct {
	app *fiber.App
}

func NewApp() *App {
	return &App{
		app: fiber.New(),
	}
}

func (a *App) Run(config Config) {
	baseCurrency, err := model.CurrencyFromString(config.BaseCurrencyStr)
	if err != nil {
		log.Fatal(err)
	}
	quoteCurrency, err := model.CurrencyFromString(config.QuoteCurrencyStr)
	if err != nil {
		log.Fatal(err)
	}

	cryptoParser := parser.NewBinanceCryptoParser(config.CryptoParserBaseURL)

	cryptoMailer := mailer.NewMailer("smtp.gmail.com", "587", config.CryptoMailerSenderEmail, config.CryptoMailerSenderPassword)

	subscriberRepository := repository.NewSubscriberFileRepository(config.SubscriberRepositoryEmailsFilePath)

	mailerService := service.NewMailerService(subscriberRepository, cryptoMailer)

	rateHandler := handler.NewRateHandler(cryptoParser, baseCurrency, quoteCurrency)

	mailerHandler := handler.NewMailerHandler(mailerService, cryptoParser, subscriberRepository, validator.New(), baseCurrency, quoteCurrency)

	api := a.app.Group("/api")

	api.Get("/rate", rateHandler.GetExchangeRate)
	api.Post("/subscribe", mailerHandler.Subscribe)
	api.Post("/sendEmails", mailerHandler.SendExchangeRate)

	if err := a.app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func (a *App) Shutdown() {
	if err := a.app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
