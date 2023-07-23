package e2e_tests

import (
	"bitcoin-exchange-rate/pkg/database"
	"bitcoin-exchange-rate/pkg/webserver"
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

	db := database.NewDB(database.DBConfig{
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DriverName: os.Getenv("DB_DRIVER_NAME"),
	})

	app := webserver.NewApp(db)
	go app.Run(webserver.Config{
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

	request, err := http.NewRequest(http.MethodGet, "http://localhost:3000/api/rate", nil)
	require.NoError(t, err)

	client := &http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var data struct {
		Rate float64 `json:"rate"`
	}

	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	assert.Greater(t, data.Rate, 0.0)

	defer func() {
		err := response.Body.Close()
		require.NoError(t, err)
		app.Shutdown()
	}()
}

func TestRateHandler_GetExchangeRate_Failure(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err)

	db := database.NewDB(database.DBConfig{
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DriverName: os.Getenv("DB_DRIVER_NAME"),
	})

	app := webserver.NewApp(db)
	go app.Run(webserver.Config{
		BinanceCryptoProviderBaseURL:  os.Getenv("BINANCE_BASE_URL"),
		CoinAPICryptoProviderBaseURL:  os.Getenv("COIN_API_BASE_URL"),
		CoinBaseCryptoProviderBaseURL: os.Getenv("COINBASE_BASE_URL"),
		CoinAPICryptoProviderKey:      "invalid-key",
		DefaultProviderName:           os.Getenv("DEFAULT_PROVIDER_NAME"),
		CryptoMailerSenderEmail:       os.Getenv("SENDER_EMAIL"),
		CryptoMailerSenderPassword:    os.Getenv("SENDER_PASSWORD"),
		BaseCurrencyStr:               os.Getenv("BASE_CURRENCY"),
		QuoteCurrencyStr:              os.Getenv("BASE_CURRENCY"),
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
