package handler

import (
	"bitcoin-exchange-rate/pkg/parser"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRateHandler_GetExchangeRateRate(t *testing.T) {

	tests := []struct {
		name               string
		baseURL            string
		expectedStatusCode int
	}{
		{
			name:               "Get rate successfully",
			baseURL:            "https://api.binance.com/api/v3/ticker",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Get rate failure (invalid url)",
			baseURL:            "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := godotenv.Load("../../.env.test")
			require.NoError(t, err, "Failed to load .env.test file")

			err = os.Setenv("BASE_URL", test.baseURL)
			require.NoError(t, err, "Failed to set BASE_URL in .env.test file")

			app := fiber.New()
			api := app.Group("/api")
			cryptoParser := parser.NewBinanceCryptoParser(os.Getenv("BASE_URL"))
			rateHandler := NewRateHandler(cryptoParser)
			api.Get("/rate", rateHandler.GetExchangeRate)

			req := httptest.NewRequest(http.MethodGet, "/api/rate", nil)
			resp, err := app.Test(req)
			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					t.Fatal(err)
				}
			}(resp.Body)

			require.NoError(t, err)
			require.Equal(t, test.expectedStatusCode, resp.StatusCode)
		})
	}
}
