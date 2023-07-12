package binance_provider

import (
	"bitcoin-exchange-rate/internal/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetExchangeRate_BinanceAPI_Success(t *testing.T) {
	err := godotenv.Load("../../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	baseCurrency, quoteCurrency := model.GetCurrencies(os.Getenv("BASE_CURRENCY"), os.Getenv("QUOTE_CURRENCY"))

	binanceProvider := NewBinanceCryptoProvider(os.Getenv("BINANCE_BASE_URL"))
	rate, err := binanceProvider.GetExchangeRateValue(baseCurrency, quoteCurrency)

	require.NoError(t, err, "Failure occurs while parsing exchange rate using BinanceApi")
	assert.Greater(t, rate, 0.0)
}

func TestGetExchangeRate_BinanceAPI_Failure(t *testing.T) {
	err := godotenv.Load("../../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	binanceProvider := NewBinanceCryptoProvider("invalid-url")

	baseCurrency, quoteCurrency := model.GetCurrencies(os.Getenv("BASE_CURRENCY"), os.Getenv("QUOTE_CURRENCY"))

	_, err = binanceProvider.GetExchangeRateValue(baseCurrency, quoteCurrency)
	assert.Error(t, err)
}
