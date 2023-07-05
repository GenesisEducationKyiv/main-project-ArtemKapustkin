package pkg

import (
	"bitcoin-exchange-rate/internal/model"
	"log"
	"reflect"
)

type loggingProvider struct {
	provider RateProvider
}

func NewLoggingProvider(provider RateProvider) *loggingProvider {
	return &loggingProvider{
		provider: provider,
	}
}

func (l *loggingProvider) GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error) {
	rate, err := l.provider.GetExchangeRateValue(baseCurrency, quoteCurrency)
	if err != nil {
		return 0, err
	}

	log.Printf("%s provides rate: %.2f", l.getProviderName(), rate)

	return rate, nil
}

func (l *loggingProvider) getProviderName() string {
	return reflect.TypeOf(l.provider).Elem().Name()
}
