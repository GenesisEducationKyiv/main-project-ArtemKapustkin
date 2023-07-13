package service

import (
	"bitcoin-exchange-rate/modules/rate_module/model"
)

type RateProvider interface {
	GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error)
}

type ExchangeRateService struct {
	rateProvider  RateProvider
	baseCurrency  model.Currency
	quoteCurrency model.Currency
}

func NewExchangeRateService(rateProvider RateProvider, baseCurrency, quoteCurrency model.Currency) *ExchangeRateService {
	return &ExchangeRateService{
		rateProvider:  rateProvider,
		baseCurrency:  baseCurrency,
		quoteCurrency: quoteCurrency,
	}
}

func (s *ExchangeRateService) GetRate() (float64, error) {
	rate, err := s.rateProvider.GetExchangeRateValue(s.baseCurrency, s.quoteCurrency)
	if err != nil {
		return 0, err
	}
	return rate, nil
}
