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
	BinanceCryptoParserBaseURL         string
	CoinCryptoParserBaseURL            string
	CoinCryptoParserAPIKey             string
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

	coinCryptoParser := parser.NewCoinCryptoParser(config.CoinCryptoParserBaseURL, config.CoinCryptoParserAPIKey)
	binanceCryptoParser := parser.NewBinanceCryptoParser(config.BinanceCryptoParserBaseURL)

	loggingCoinCryptoParser := parser.NewLoggingProvider(coinCryptoParser)
	loggingBinanceCryptoParser := parser.NewLoggingProvider(binanceCryptoParser)

	coinProviderNode := parser.NewRateProviderNode(loggingCoinCryptoParser)
	binanceProviderNode := parser.NewRateProviderNode(loggingBinanceCryptoParser)

	binanceProviderNode.SetNext(coinProviderNode)

	cryptoMailer := mailer.NewMailer("smtp.gmail.com", "587", config.CryptoMailerSenderEmail, config.CryptoMailerSenderPassword)

	subscriberRepository := repository.NewSubscriberFileRepository(config.SubscriberRepositoryEmailsFilePath)

	mailerService := service.NewMailerService(subscriberRepository, cryptoMailer)

	rateHandler := handler.NewRateHandler(binanceProviderNode, baseCurrency, quoteCurrency)

	mailerHandler := handler.NewMailerHandler(mailerService, binanceProviderNode, subscriberRepository, validator.New(), baseCurrency, quoteCurrency)

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
