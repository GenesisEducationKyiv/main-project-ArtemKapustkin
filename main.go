package main

import (
	"bitcoin-exchange-rate/internal/webserver"
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
		BinanceCryptoParserBaseURL:         os.Getenv("BINANCE_BASE_URL"),
		CoinCryptoParserBaseURL:            os.Getenv("COIN_BASE_URL"),
		CoinCryptoParserAPIKey:             os.Getenv("COIN_API_KEY"),
		CryptoMailerSenderEmail:            os.Getenv("SENDER_EMAIL"),
		CryptoMailerSenderPassword:         os.Getenv("SENDER_PASSWORD"),
		SubscriberRepositoryEmailsFilePath: os.Getenv("EMAILS_FILEPATH"),
		BaseCurrencyStr:                    os.Getenv("BASE_CURRENCY"),
		QuoteCurrencyStr:                   os.Getenv("QUOTE_CURRENCY"),
	})

	defer app.Shutdown()
}
