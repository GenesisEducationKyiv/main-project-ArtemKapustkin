package coinbase_provider

import (
	"bitcoin-exchange-rate/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type CoinBaseAPICryptoProvider struct {
	baseURL string
}

func NewCoinBaseAPICryptoProvider(baseURL string) *CoinBaseAPICryptoProvider {
	return &CoinBaseAPICryptoProvider{
		baseURL: baseURL,
	}
}

type coinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

func (p *CoinBaseAPICryptoProvider) GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error) {
	client := http.Client{}
	requestURL := fmt.Sprintf("%s/prices/%s-%s/spot", p.baseURL, baseCurrency, quoteCurrency)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result coinbaseResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		return 0, err
	}

	return price, err
}
