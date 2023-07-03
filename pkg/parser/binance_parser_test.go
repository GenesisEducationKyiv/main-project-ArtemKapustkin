package parser

import (
	"bitcoin-exchange-rate/internal/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func GetCurrenciesFromEnvFile() (model.Currency, model.Currency) {
	baseCurrency, err := model.CurrencyFromString(os.Getenv("BASE_CURRENCY"))
	if err != nil {
		log.Fatal(err)
	}

	quoteCurrency, err := model.CurrencyFromString(os.Getenv("QUOTE_CURRENCY"))
	if err != nil {
		log.Fatal(err)
	}

	return baseCurrency, quoteCurrency
}

func TestGetExchangeRate_BinanceAPI_Success(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	baseCurrency, quoteCurrency := GetCurrenciesFromEnvFile()

	binanceParser := NewBinanceCryptoParser(os.Getenv("BASE_URL"))
	rate, err := binanceParser.GetExchangeRateValue(baseCurrency, quoteCurrency)

	require.NoError(t, err, "Failure occurs while parsing exchange rate using BinanceApi")
	assert.Greater(t, rate, 0)
}

func TestGetExchangeRate_BinanceAPI_Failure(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	binanceParser := NewBinanceCryptoParser("invalid-url")

	baseCurrency, quoteCurrency := GetCurrenciesFromEnvFile()

	_, err = binanceParser.GetExchangeRateValue(baseCurrency, quoteCurrency)
	assert.Error(t, err)
}
