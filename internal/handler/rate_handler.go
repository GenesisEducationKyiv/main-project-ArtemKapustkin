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
	exchangeRateProvider ExchangeRateClient

	exchangeRateBaseCurrency  model.Currency
	exchangeRateQuoteCurrency model.Currency
}

func NewRateHandler(
	exchangeRateProvider ExchangeRateClient,
	baseCurrency model.Currency,
	quoteCurrency model.Currency,
) *RateHandler {
	return &RateHandler{
		exchangeRateProvider:      exchangeRateProvider,
		exchangeRateBaseCurrency:  baseCurrency,
		exchangeRateQuoteCurrency: quoteCurrency,
	}
}

func (h *RateHandler) GetExchangeRate(c *fiber.Ctx) error {
	rate, err := h.exchangeRateProvider.GetExchangeRateValue(h.exchangeRateBaseCurrency, h.exchangeRateQuoteCurrency)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.JSON(rate)
}
