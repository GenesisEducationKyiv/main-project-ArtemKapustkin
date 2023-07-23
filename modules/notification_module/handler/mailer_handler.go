package handler

import (
	"bitcoin-exchange-rate/modules/notification_module/model"
	"bitcoin-exchange-rate/pkg/presenter"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type MailerService interface {
	SendExchangeRate() error
	Subscribe(subscriber *model.Subscriber) error
}

type MailerHandler struct {
	mailerService MailerService
	validator     *validator.Validate
	presenter     presenter.ResponsePresenter
}

func NewMailerHandler(
	mailerService MailerService,
	validator *validator.Validate,
	presenter presenter.ResponsePresenter,
) *MailerHandler {
	return &MailerHandler{
		mailerService: mailerService,
		validator:     validator,
		presenter:     presenter,
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

	err := h.mailerService.Subscribe(model.NewSubscriber(payload.Email))
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
