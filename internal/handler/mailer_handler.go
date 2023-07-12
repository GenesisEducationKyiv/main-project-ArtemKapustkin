package handler

import (
	"bitcoin-exchange-rate/internal/model"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type SubscriptionRepository interface {
	Create(subscriber *model.Subscriber) error
}

type MailerService interface {
	SendExchangeRate() error
}

type MailerHandler struct {
	mailerService          MailerService
	subscriptionRepository SubscriptionRepository
	validator              *validator.Validate
	presenter              ResponsePresenter
}

func NewMailerHandler(
	mailerService MailerService,
	subscriptionRepository SubscriptionRepository,
	validator *validator.Validate,
	presenter ResponsePresenter,
) *MailerHandler {
	return &MailerHandler{
		mailerService:          mailerService,
		subscriptionRepository: subscriptionRepository,
		validator:              validator,
		presenter:              presenter,
	}
}

func (h *MailerHandler) SendExchangeRate(c *fiber.Ctx) error {
	if err := h.mailerService.SendExchangeRate(); err != nil {
		if errors.Is(err, model.ErrSubscriberFileIsEmpty) {
			return h.presenter.PresentError(c, http.StatusBadRequest, err)
		}

		return h.presenter.PresentError(c, http.StatusInternalServerError, err)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *MailerHandler) Subscribe(c *fiber.Ctx) error {
	var payload subscribeDTO

	if err := c.BodyParser(&payload); err != nil {
		return h.presenter.PresentError(c, http.StatusBadRequest, err)
	}

	if err := h.validator.Struct(&payload); err != nil {
		return h.presenter.PresentError(c, http.StatusBadRequest, err)
	}

	err := h.subscriptionRepository.Create(model.NewSubscriber(payload.Email))
	if err != nil {
		if errors.Is(err, model.ErrSubscriberAlreadyExist) {
			return h.presenter.PresentError(c, http.StatusConflict, err)
		}
		return h.presenter.PresentError(c, http.StatusInternalServerError, err)
	}

	return c.SendStatus(http.StatusOK)
}

type subscribeDTO struct {
	Email string `validate:"required,email"`
}
