package webserver

import (
	"bitcoin-exchange-rate/internal/handler"
	"bitcoin-exchange-rate/internal/model"
	"bitcoin-exchange-rate/internal/repository"
	"bitcoin-exchange-rate/internal/service"
	"bitcoin-exchange-rate/pkg"
	"bitcoin-exchange-rate/pkg/mailer"
	"bitcoin-exchange-rate/pkg/rate_providers/binance_provider"
	"bitcoin-exchange-rate/pkg/rate_providers/coinapi_provider"
	"bitcoin-exchange-rate/pkg/rate_providers/coinbase_provider"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Config struct {
	BinanceCryptoProviderBaseURL       string
	CoinAPICryptoProviderBaseURL       string
	CoinBaseCryptoProviderBaseURL      string
	CoinAPICryptoProviderKey           string
	DefaultProviderName                string
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

func getConfiguredProvider(config Config) pkg.RateProvider {
	binanceCryptoProvider := binance_provider.NewBinanceCryptoProvider(config.BinanceCryptoProviderBaseURL)
	coinAPICryptoProvider := coinapi_provider.NewCoinAPICryptoProvider(config.CoinAPICryptoProviderBaseURL, config.CoinAPICryptoProviderKey)
	coinBaseCryptoProvider := coinbase_provider.NewCoinBaseAPICryptoProvider(config.CoinBaseCryptoProviderBaseURL)

	loggingBinanceCryptoProvider := pkg.NewLoggingProvider(binanceCryptoProvider)
	loggingCoinAPICryptoProvider := pkg.NewLoggingProvider(coinAPICryptoProvider)
	loggingCoinBaseCryptoProvider := pkg.NewLoggingProvider(coinBaseCryptoProvider)

	binanceProviderNode := pkg.NewRateProviderNode(loggingBinanceCryptoProvider)
	coinAPIProviderNode := pkg.NewRateProviderNode(loggingCoinAPICryptoProvider)
	coinBaseProviderNode := pkg.NewRateProviderNode(loggingCoinBaseCryptoProvider)

	providersChain := pkg.NewProvidersChain()

	err := providersChain.RegisterProvider("binance", binanceProviderNode, coinAPIProviderNode)
	if err != nil {
		return nil
	}
	err = providersChain.RegisterProvider("coinapi", coinAPIProviderNode, coinBaseProviderNode)
	if err != nil {
		return nil
	}
	err = providersChain.RegisterProvider("coinbase", coinBaseProviderNode, nil)
	if err != nil {
		return nil
	}

	return providersChain.GetProvider(config.DefaultProviderName)
}

func (a *App) Run(config Config) {
	baseCurrency, quoteCurrency := model.GetCurrencies(config.BaseCurrencyStr, config.QuoteCurrencyStr)

	rateProvider := getConfiguredProvider(config)

	cryptoMailer := mailer.NewMailer("smtp.gmail.com", "587", config.CryptoMailerSenderEmail, config.CryptoMailerSenderPassword)

	subscriberRepository := repository.NewSubscriberFileRepository(config.SubscriberRepositoryEmailsFilePath)

	mailerService := service.NewMailerService(subscriberRepository, cryptoMailer)

	rateHandler := handler.NewRateHandler(rateProvider, baseCurrency, quoteCurrency)

	mailerHandler := handler.NewMailerHandler(mailerService, rateProvider, subscriberRepository, validator.New(), baseCurrency, quoteCurrency)

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
