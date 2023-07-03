package parser

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetExchangeRate_CoinAPI_Success(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	baseCurrency, quoteCurrency := GetCurrenciesFromEnvFile()

	coinParser := NewCoinCryptoParser(os.Getenv("COIN_BASE_URL"), os.Getenv("COIN_API_KEY"))
	rate, err := coinParser.GetExchangeRateValue(baseCurrency, quoteCurrency)

	require.NoError(t, err, "Failure occurs while parsing exchange rate using CoinAPI")
	assert.Greater(t, rate, 0)
}

func TestGetExchangeRate_CoinAPI_Failure(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	coinParser := NewCoinCryptoParser("invalid-url", os.Getenv("COIN_API_KEY"))

	baseCurrency, quoteCurrency := GetCurrenciesFromEnvFile()

	_, err = coinParser.GetExchangeRateValue(baseCurrency, quoteCurrency)
	assert.Error(t, err)
}
