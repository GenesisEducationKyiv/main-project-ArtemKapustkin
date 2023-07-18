package service

import (
	"bitcoin-exchange-rate/modules/rate_module/model"
)

type RateProvider interface {
	GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error)
}

type Logger interface {
	Error(message string)
	Info(message string)
}

type ExchangeRateService struct {
	rateProvider  RateProvider
	baseCurrency  model.Currency
	quoteCurrency model.Currency
	logger        Logger
}

func NewExchangeRateService(rateProvider RateProvider, baseCurrency, quoteCurrency model.Currency, logger Logger) *ExchangeRateService {
	return &ExchangeRateService{
		rateProvider:  rateProvider,
		baseCurrency:  baseCurrency,
		quoteCurrency: quoteCurrency,
		logger:        logger,
	}
}

func (s *ExchangeRateService) GetRate() (float64, error) {
	rate, err := s.rateProvider.GetExchangeRateValue(s.baseCurrency, s.quoteCurrency)
	if err != nil {
		s.logger.Error(err.Error())
		return 0, err
	}
	return rate, nil
}
