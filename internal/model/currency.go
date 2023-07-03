package model

import (
	"errors"
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
