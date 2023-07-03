package handler

import (
	"bitcoin-exchange-rate/internal/model"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ExchangeRateClient interface {
	GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error)
}

type RateHandler struct {
	exchangeRateParser ExchangeRateClient

	exchangeRateBaseCurrency  model.Currency
	exchangeRateQuoteCurrency model.Currency
}

func NewRateHandler(
	exchangeRateParser ExchangeRateClient,
	baseCurrency model.Currency,
	quoteCurrency model.Currency,
) *RateHandler {
	return &RateHandler{
		exchangeRateParser:        exchangeRateParser,
		exchangeRateBaseCurrency:  baseCurrency,
		exchangeRateQuoteCurrency: quoteCurrency,
	}
}

func (h *RateHandler) GetExchangeRate(c *fiber.Ctx) error {
	rate, err := h.exchangeRateParser.GetExchangeRateValue(h.exchangeRateBaseCurrency, h.exchangeRateQuoteCurrency)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.JSON(rate)
}
