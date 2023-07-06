package coinbase_provider

import (
	"bitcoin-exchange-rate/internal/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestGetExchangeRate_CoinBaseAPI_Success(t *testing.T) {
	err := godotenv.Load("../../../.env.test")
	require.NoError(t, err, "Failed to load .env.test file")

	baseCurrency, quoteCurrency := model.GetCurrencies(os.Getenv("BASE_CURRENCY"), os.Getenv("QUOTE_CURRENCY"))

	coinBaseAPIProvider := NewCoinBaseAPICryptoProvider(os.Getenv("COIN_API_BASE_URL"))
	rate, err := coinBaseAPIProvider.GetExchangeRateValue(baseCurrency, quoteCurrency)
	log.Printf("Currency Rate: %.2f", rate)
	require.NoError(t, err, "Failure occurs while parsing exchange rate using CoinAPI")
	assert.Greater(t, rate, 0.0)
}
