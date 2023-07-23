package main

import (
	"bitcoin-exchange-rate/pkg/database"
	"bitcoin-exchange-rate/pkg/webserver"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.NewDB(database.DBConfig{
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DriverName: os.Getenv("DB_DRIVER_NAME"),
	})

	app := webserver.NewApp(db)

	app.Run(webserver.Config{
		BinanceCryptoProviderBaseURL:  os.Getenv("BINANCE_BASE_URL"),
		CoinAPICryptoProviderBaseURL:  os.Getenv("COIN_API_BASE_URL"),
		CoinBaseCryptoProviderBaseURL: os.Getenv("COINBASE_BASE_URL"),
		CoinAPICryptoProviderKey:      os.Getenv("COIN_API_KEY"),
		DefaultProviderName:           os.Getenv("DEFAULT_PROVIDER_NAME"),
		CryptoMailerSenderEmail:       os.Getenv("SENDER_EMAIL"),
		CryptoMailerSenderPassword:    os.Getenv("SENDER_PASSWORD"),
		BaseCurrencyStr:               os.Getenv("BASE_CURRENCY"),
		QuoteCurrencyStr:              os.Getenv("QUOTE_CURRENCY"),
	})

	defer app.Shutdown()
}
