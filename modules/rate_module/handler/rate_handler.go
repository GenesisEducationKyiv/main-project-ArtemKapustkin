package handler

import (
	"bitcoin-exchange-rate/pkg/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

var ErrProviderGetRate = errors.New("there is an error, while parsing rate")

type ExchangeRateService interface {
	GetRate() (float64, error)
}

type RateHandler struct {
	exchangeRateService ExchangeRateService
	presenter           presenter.ResponsePresenter
}

func NewRateHandler(
	exchangeRateService ExchangeRateService,
	presenter presenter.ResponsePresenter,
) *RateHandler {
	return &RateHandler{
		exchangeRateService: exchangeRateService,
		presenter:           presenter,
	}
}

func (h *RateHandler) GetExchangeRate(c *fiber.Ctx) error {
	rate, err := h.exchangeRateService.GetRate()
	if err != nil || rate == 0 {
		return h.presenter.PresentError(c, http.StatusBadRequest, ErrProviderGetRate)
	}

	return h.presenter.PresentExchangeRate(c, rate)
}
