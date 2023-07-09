package handler

import (
	"bitcoin-exchange-rate/internal/model"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

type ExchangeRateProvider interface {
	GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error)
}

type RateHandler struct {
	exchangeRateProvider      ExchangeRateProvider
	presenter                 ResponsePresenter
	exchangeRateBaseCurrency  model.Currency
	exchangeRateQuoteCurrency model.Currency
}

func NewRateHandler(
	exchangeRateProvider ExchangeRateProvider,
	presenter ResponsePresenter,
	baseCurrency model.Currency,
	quoteCurrency model.Currency,
) *RateHandler {
	return &RateHandler{
		exchangeRateProvider:      exchangeRateProvider,
		presenter:                 presenter,
		exchangeRateBaseCurrency:  baseCurrency,
		exchangeRateQuoteCurrency: quoteCurrency,
	}
}

func (h *RateHandler) GetExchangeRate(c *fiber.Ctx) error {
	rate, err := h.exchangeRateProvider.GetExchangeRateValue(h.exchangeRateBaseCurrency, h.exchangeRateQuoteCurrency)
	if err != nil || rate == 0 {
		log.Println("error: ", err)
		return h.presenter.PresentError(c, http.StatusBadRequest, errors.New("blablabla"))
	}

	return c.JSON(rate)
}
