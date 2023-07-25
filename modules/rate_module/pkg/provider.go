package pkg

import (
	"bitcoin-exchange-rate/modules/rate_module/model"
)

type RateProvider interface {
	GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error)
}

type RateProviderNode struct {
	provider RateProvider
	next     RateProvider
}

func NewRateProviderNode(provider RateProvider) *RateProviderNode {
	return &RateProviderNode{
		provider: provider,
	}
}

func (c *RateProviderNode) GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error) {
	rate, err := c.provider.GetExchangeRateValue(baseCurrency, quoteCurrency)
	if err != nil && c.next != nil {
		return c.next.GetExchangeRateValue(baseCurrency, quoteCurrency)
	}

	return rate, nil
}

func (c *RateProviderNode) SetNext(provider ProviderNode) {
	c.next = provider
}
