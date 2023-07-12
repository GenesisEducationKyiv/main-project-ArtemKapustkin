package model

import (
	"errors"
	"log"
)

type Currency string

const (
	BTC Currency = "BTC"
	UAH Currency = "UAH"
)

func CurrencyFromString(currencyStr string) (Currency, error) {
	switch currencyStr {
	case "BTC":
		return BTC, nil
	case "UAH":
		return UAH, nil
	default:
		return "", errors.New("unknown currency")
	}
}

func GetCurrencies(baseCurrencyStr, quoteCurrencyStr string) (Currency, Currency) {
	baseCurrency, err := CurrencyFromString(baseCurrencyStr)
	if err != nil {
		log.Fatal(err)
	}

	quoteCurrency, err := CurrencyFromString(quoteCurrencyStr)
	if err != nil {
		log.Fatal(err)
	}

	return baseCurrency, quoteCurrency
}
