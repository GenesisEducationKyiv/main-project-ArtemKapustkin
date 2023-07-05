package coinapi_provider

import (
	"bitcoin-exchange-rate/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CoinAPICryptoProvider struct {
	baseURL string
	apiKey  string
}

func NewCoinAPICryptoProvider(baseURL, apiKey string) *CoinAPICryptoProvider {
	return &CoinAPICryptoProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

type ExchangeRate struct {
	Time         string  `json:"time"`
	AssetIDBase  string  `json:"asset_id_base"`
	AssetIDQuote string  `json:"asset_id_quote"`
	Rate         float64 `json:"rate"`
}

func (p *CoinAPICryptoProvider) GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error) {
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/%s/%s", p.baseURL, baseCurrency, quoteCurrency)

	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return 0, err
	}
	request.Header.Set("X-CoinAPI-Key", p.apiKey)

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("error closing response body: %s", err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var exchangeRate ExchangeRate
	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		log.Println("failed to parse JSON: ", err)
		return 0, err
	}
	return exchangeRate.Rate, nil
}
