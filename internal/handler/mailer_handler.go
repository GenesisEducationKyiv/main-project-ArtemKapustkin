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
	mailerService          MailerService
	exchangeRateService    ExchangeRateService
	subscriptionRepository SubscriptionRepository
	validator              *validator.Validate
	presenter              ResponsePresenter
}

func NewMailerHandler(
	mailerService MailerService,
	exchangeRateService ExchangeRateService,
	subscriptionRepository SubscriptionRepository,
	validator *validator.Validate,
	presenter ResponsePresenter,
) *MailerHandler {
	return &MailerHandler{
		mailerService:          mailerService,
		exchangeRateService:    exchangeRateService,
		subscriptionRepository: subscriptionRepository,
		validator:              validator,
		presenter:              presenter,
	}
}

func (h *MailerHandler) SendExchangeRate(c *fiber.Ctx) error {
	value, err := h.exchangeRateService.GetRate()
	if err != nil {
		return h.presenter.PresentError(c, http.StatusInternalServerError, err)
	}

	err = h.mailerService.SendValueToAllEmails(model.NewEmailMessage(strconv.FormatFloat(value, 'f', 2, 64)))
	if err != nil {
		return h.presenter.PresentError(c, http.StatusBadRequest, err)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *MailerHandler) Subscribe(c *fiber.Ctx) error {
	var payload subscribeDTO

	if err := c.BodyParser(&payload); err != nil {
		return h.presenter.PresentError(c, http.StatusBadRequest, err)
	}

	if h.validator.Struct(&payload) != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	err := h.subscriptionRepository.Create(model.NewSubscriber(payload.Email))
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExist) {
			return h.presenter.PresentError(c, http.StatusConflict, err)
		}
		return h.presenter.PresentError(c, http.StatusInternalServerError, err)
	}

	return c.SendStatus(http.StatusOK)
}

type subscribeDTO struct {
	Email string `validate:"required,email"`
}
