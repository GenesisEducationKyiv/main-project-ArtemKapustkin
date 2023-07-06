package coinapi_provider

import (
	"bitcoin-exchange-rate/internal/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetExchangeRate_CoinAPI_Success(t *testing.T) {
	err := godotenv.Load("../../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	baseCurrency, quoteCurrency := model.GetCurrencies(os.Getenv("BASE_CURRENCY"), os.Getenv("QUOTE_CURRENCY"))

	coinAPIProvider := NewCoinAPICryptoProvider(os.Getenv("COIN_API_BASE_URL"), os.Getenv("COIN_API_KEY"))
	rate, err := coinAPIProvider.GetExchangeRateValue(baseCurrency, quoteCurrency)

	require.NoError(t, err, "Failure occurs while parsing exchange rate using CoinAPI")
	assert.Greater(t, rate, 0)
}

func TestGetExchangeRate_CoinAPI_Failure(t *testing.T) {
	err := godotenv.Load("../../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	coinAPIProvider := NewCoinAPICryptoProvider("invalid-url", os.Getenv("COIN_API_KEY"))

	baseCurrency, quoteCurrency := model.GetCurrencies(os.Getenv("BASE_CURRENCY"), os.Getenv("QUOTE_CURRENCY"))

	_, err = coinAPIProvider.GetExchangeRateValue(baseCurrency, quoteCurrency)
	assert.Error(t, err)
}
