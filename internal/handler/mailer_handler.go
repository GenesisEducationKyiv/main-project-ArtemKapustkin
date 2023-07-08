package handler

import (
	"bitcoin-exchange-rate/internal/model"
	"bitcoin-exchange-rate/internal/repository"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SubscriptionRepository interface {
	Create(subscriber *model.Subscriber) error
}

type MailerService interface {
	SendValueToAllEmails(message model.EmailMessage) error
}

type MailerHandler struct {
	mailerService             MailerService
	exchangeRateService       ExchangeRateProvider
	subscriptionRepository    SubscriptionRepository
	validator                 *validator.Validate
	presenter                 ResponsePresenter
	exchangeRateBaseCurrency  model.Currency
	exchangeRateQuoteCurrency model.Currency
}

func NewMailerHandler(
	mailerService MailerService,
	exchangeRateService ExchangeRateProvider,
	subscriptionRepository SubscriptionRepository,
	validator *validator.Validate,
	presenter ResponsePresenter,
	baseCurrency model.Currency,
	quoteCurrency model.Currency,
) *MailerHandler {
	return &MailerHandler{
		mailerService:             mailerService,
		exchangeRateService:       exchangeRateService,
		subscriptionRepository:    subscriptionRepository,
		validator:                 validator,
		presenter:                 presenter,
		exchangeRateBaseCurrency:  baseCurrency,
		exchangeRateQuoteCurrency: quoteCurrency,
	}
}

func (h *MailerHandler) SendExchangeRate(c *fiber.Ctx) error {
	value, err := h.exchangeRateService.GetExchangeRateValue(h.exchangeRateBaseCurrency, h.exchangeRateQuoteCurrency)
	if err != nil {
		return h.presenter.PresentError(c.Status(http.StatusInternalServerError), err)
	}

	err = h.mailerService.SendValueToAllEmails(model.NewEmailMessage(strconv.FormatFloat(value, 'f', 2, 64)))
	if err != nil {
		return h.presenter.PresentError(c.Status(http.StatusBadRequest), err)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *MailerHandler) Subscribe(c *fiber.Ctx) error {
	var payload subscribeDTO

	if err := c.BodyParser(&payload); err != nil {
		return h.presenter.PresentError(c.Status(http.StatusBadRequest), err)
	}

	if h.validator.Struct(&payload) != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	err := h.subscriptionRepository.Create(model.NewSubscriber(payload.Email))
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExist) {
			return h.presenter.PresentError(c.Status(http.StatusConflict), err)
		}
		return h.presenter.PresentError(c.Status(http.StatusInternalServerError), err)
	}

	return c.SendStatus(http.StatusOK)
}

type subscribeDTO struct {
	Email string `validate:"required,email"`
}
