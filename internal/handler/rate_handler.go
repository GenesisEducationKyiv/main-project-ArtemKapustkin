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
	binanceParser ExchangeRateClient
}

func NewRateHandler(binanceParser ExchangeRateClient) *RateHandler {
	return &RateHandler{
		binanceParser: binanceParser,
	}
}

func (h *RateHandler) GetExchangeRate(c *fiber.Ctx) error {
	rate, err := h.binanceParser.GetExchangeRateValue("BTC", "UAH")
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.JSON(rate)
}
