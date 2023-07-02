package parser

import (
	"bitcoin-exchange-rate/internal/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	binanceParser := NewBinanceCryptoParser(os.Getenv("BASE_URL"))
	rate, err := binanceParser.GetExchangeRateValue(model.BTC, model.UAH)

	require.NoError(t, err, "Failure occurs while parsing exchange rate")
	assert.Greater(t, rate, 0.0)
}

func TestGetExchangeRateFault(t *testing.T) {
	binanceParser := NewBinanceCryptoParser("invalid-url")
	_, err := binanceParser.GetExchangeRateValue(model.BTC, model.UAH)
	assert.Error(t, err)
}
