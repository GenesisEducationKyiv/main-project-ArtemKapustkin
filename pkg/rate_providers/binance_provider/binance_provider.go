package binance_provider

import (
	"bitcoin-exchange-rate/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type BinanceCryptoProvider struct {
	baseURL string
}

func NewBinanceCryptoProvider(baseURL string) *BinanceCryptoProvider {
	return &BinanceCryptoProvider{
		baseURL: baseURL,
	}
}

func (p *BinanceCryptoProvider) GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error) {
	requestURL := fmt.Sprintf("%s/price?symbol=%s%s", p.baseURL, baseCurrency, quoteCurrency)

	response, err := http.Get(requestURL)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("error closing response body: %s", err)
		}
	}()

	var result struct {
		Price string `json:"price"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0, err
	}

	rate, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, err
	}

	return rate, nil
}
