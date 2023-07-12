package pkg

import "bitcoin-exchange-rate/internal/model"

type RateProvider interface {
	GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error)
}

type RateProviderNode interface {
	RateProvider
	SetNext(provider RateProviderNode)
}

type rateProviderNode struct {
	provider RateProvider
	next     RateProvider
}

func NewRateProviderNode(provider RateProvider) *rateProviderNode {
	return &rateProviderNode{
		provider: provider,
	}
}

func (c *rateProviderNode) GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error) {
	rate, err := c.provider.GetExchangeRateValue(baseCurrency, quoteCurrency)
	if err != nil && c.next != nil {
		return c.next.GetExchangeRateValue(baseCurrency, quoteCurrency)
	}

	return rate, nil
}

func (c *rateProviderNode) SetNext(provider RateProviderNode) {
	c.next = provider
}
