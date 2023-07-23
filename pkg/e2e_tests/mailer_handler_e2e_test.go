package e2e_tests

import (
	"bitcoin-exchange-rate/pkg/database"
	"bitcoin-exchange-rate/pkg/webserver"
	"bytes"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestMailerHandler_Subscribe(t *testing.T) {
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

	tests := []struct {
		name               string
		expectedStatusCode int
		body               string
	}{
		{
			name:               "Subscribe successful",
			expectedStatusCode: http.StatusOK,
			body:               `{"email": "example@example.com"}`,
		},
		{
			name:               "Invalid request body",
			expectedStatusCode: http.StatusBadRequest,
			body:               ``,
		},
		{
			name:               "Invalid email address",
			expectedStatusCode: http.StatusBadRequest,
			body:               `{"email": "invalid-email"}`,
		},
		{
			name:               "Already subscribed",
			expectedStatusCode: http.StatusConflict,
			body:               `{"email": "example@example.com"}`,
		},
	}

	client := &http.Client{}

	for _, test := range tests {
		request, err := http.NewRequest(http.MethodPost, "http://localhost:3000/api/subscribe", bytes.NewBufferString(test.body))
		require.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		require.NoError(t, err)

		err = response.Body.Close()
		require.NoError(t, err)

		require.Equal(t, test.expectedStatusCode, response.StatusCode)
	}

	require.NoError(t, err)

	defer app.Shutdown()
}

func TestMailerHandler_SendEmails_Success(t *testing.T) {
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

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPost, "http://localhost:3000/api/sendEmails", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	defer func() {
		err := response.Body.Close()
		require.NoError(t, err)
		app.Shutdown()
	}()
}

func TestMailerHandler_SendEmails_Failure(t *testing.T) {
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

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPost, "http://localhost:3000/api/sendEmails", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	defer func() {
		err := response.Body.Close()
		require.NoError(t, err)
		app.Shutdown()
	}()
}
