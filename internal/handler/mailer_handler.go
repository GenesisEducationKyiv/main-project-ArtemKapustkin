package handler

import (
	"bitcoin-exchange-rate/internal/model"
	"bitcoin-exchange-rate/internal/repository"
	"bitcoin-exchange-rate/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SubscriberRepository interface {
	Create(subscriber *model.Subscriber) error
}

type ParserService interface {
	GetExchangeRateValue(baseCurrency model.Currency, quoteCurrency model.Currency) (float64, error)
}

type MailerHandler struct {
	mailerService        *service.MailerService
	parserService        ParserService
	subscriberRepository SubscriberRepository
	validator            *validator.Validate

	exchangeRateBaseCurrency  model.Currency
	exchangeRateQuoteCurrency model.Currency
}

func NewMailerHandler(
	mailerService *service.MailerService,
	parserService ParserService,
	subscriberRepository SubscriberRepository,
	validator *validator.Validate,
	baseCurrency model.Currency,
	quoteCurrency model.Currency,
) *MailerHandler {
	return &MailerHandler{
		mailerService:             mailerService,
		parserService:             parserService,
		subscriberRepository:      subscriberRepository,
		validator:                 validator,
		exchangeRateBaseCurrency:  baseCurrency,
		exchangeRateQuoteCurrency: quoteCurrency,
	}
}

func (h *MailerHandler) SendExchangeRate(c *fiber.Ctx) error {
	value, err := h.parserService.GetExchangeRateValue(h.exchangeRateBaseCurrency, h.exchangeRateQuoteCurrency)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	err = h.mailerService.SendValueToAllEmails(strconv.FormatFloat(value, 'f', 2, 64))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *MailerHandler) Subscribe(c *fiber.Ctx) error {
	var payload subscribeDTO

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if h.validator.Struct(&payload) != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	err := h.subscriberRepository.Create(model.NewSubscriber(payload.Email))
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExist) {
			return c.SendStatus(http.StatusConflict)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

type subscribeDTO struct {
	Email string `validate:"required,email"`
}
