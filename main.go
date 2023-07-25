package main

import (
	"bitcoin-exchange-rate/pkg/webserver"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	app := webserver.NewApp()

	app.Run(webserver.Config{
		BinanceCryptoProviderBaseURL:       os.Getenv("BINANCE_BASE_URL"),
		CoinAPICryptoProviderBaseURL:       os.Getenv("COIN_API_BASE_URL"),
		CoinBaseCryptoProviderBaseURL:      os.Getenv("COINBASE_BASE_URL"),
		CoinAPICryptoProviderKey:           os.Getenv("COIN_API_KEY"),
		DefaultProviderName:                os.Getenv("DEFAULT_PROVIDER_NAME"),
		CryptoMailerSenderEmail:            os.Getenv("SENDER_EMAIL"),
		CryptoMailerSenderPassword:         os.Getenv("SENDER_PASSWORD"),
		SubscriberRepositoryEmailsFilePath: os.Getenv("EMAILS_FILE_PATH"),
		BaseCurrencyStr:                    os.Getenv("BASE_CURRENCY"),
		QuoteCurrencyStr:                   os.Getenv("QUOTE_CURRENCY"),
	})

	defer app.Shutdown()
}
