package handler

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ExchangeRateClient interface {
	GetExchangeRate(baseCurrency string, quoteCurrency string) (float64, error)
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
	rate, err := h.binanceParser.GetExchangeRate("BTC", "UAH")
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.JSON(rate)
}
