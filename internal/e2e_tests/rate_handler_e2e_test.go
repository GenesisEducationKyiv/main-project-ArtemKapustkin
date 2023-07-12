package e2e_tests

import (
	"bitcoin-exchange-rate/internal/webserver"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestRateHandler_GetExchangeRate_Success(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err)

	app := webserver.NewApp()
	go app.Run(webserver.Config{
		BinanceCryptoProviderBaseURL:       os.Getenv("BINANCE_BASE_URL"),
		CoinAPICryptoProviderBaseURL:       os.Getenv("COIN_API_BASE_URL"),
		CoinBaseCryptoProviderBaseURL:      os.Getenv("COINBASE_BASE_URL"),
		CoinAPICryptoProviderKey:           os.Getenv("COIN_API_KEY"),
		DefaultProviderName:                os.Getenv("DEFAULT_PROVIDER_NAME"),
		CryptoMailerSenderEmail:            os.Getenv("SENDER_EMAIL"),
		CryptoMailerSenderPassword:         os.Getenv("SENDER_PASSWORD"),
		SubscriberRepositoryEmailsFilePath: os.Getenv("TEST_FILE_PATH"),
		BaseCurrencyStr:                    os.Getenv("BASE_CURRENCY"),
		QuoteCurrencyStr:                   os.Getenv("QUOTE_CURRENCY"),
	})

	request, err := http.NewRequest(http.MethodGet, "http://localhost:3000/api/rate", nil)
	require.NoError(t, err)

	client := &http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var rate float64
	err = json.Unmarshal(body, &rate)
	require.NoError(t, err)

	assert.Greater(t, rate, 0.0)

	defer func() {
		err := response.Body.Close()
		require.NoError(t, err)
		app.Shutdown()
	}()
}

func TestRateHandler_GetExchangeRate_Failure(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err)

	app := webserver.NewApp()
	go app.Run(webserver.Config{
		BinanceCryptoProviderBaseURL:       os.Getenv("BINANCE_BASE_URL"),
		CoinAPICryptoProviderBaseURL:       os.Getenv("COIN_API_BASE_URL"),
		CoinBaseCryptoProviderBaseURL:      os.Getenv("COINBASE_BASE_URL"),
		CoinAPICryptoProviderKey:           "invalid-key",
		DefaultProviderName:                os.Getenv("DEFAULT_PROVIDER_NAME"),
		CryptoMailerSenderEmail:            os.Getenv("SENDER_EMAIL"),
		CryptoMailerSenderPassword:         os.Getenv("SENDER_PASSWORD"),
		SubscriberRepositoryEmailsFilePath: os.Getenv("TEST_FILE_PATH"),
		BaseCurrencyStr:                    os.Getenv("BASE_CURRENCY"),
		QuoteCurrencyStr:                   os.Getenv("BASE_CURRENCY"),
	})

	request, err := http.NewRequest(http.MethodGet, "http://localhost:3000/api/rate", nil)
	require.NoError(t, err)

	client := &http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	defer func() {
		err := response.Body.Close()
		require.NoError(t, err)
		app.Shutdown()
	}()
}
