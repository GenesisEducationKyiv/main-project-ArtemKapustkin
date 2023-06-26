package handler

import (
	"bitcoin-exchange-rate/internal/model"
	"bitcoin-exchange-rate/internal/repository"
	"bitcoin-exchange-rate/internal/service"
	"bitcoin-exchange-rate/pkg/parser"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SubscriberRepository interface {
	Create(subscriber *model.Subscriber) error
}

type MailerHandler struct {
	mailerService        *service.MailerService
	binanceCryptoParser  *parser.BinanceCryptoParser
	subscriberRepository SubscriberRepository
	validator            *validator.Validate

	exchangeRateBaseCurrency  model.Currency
	exchangeRateQuoteCurrency model.Currency
}

func NewMailerHandler(
	mailerService *service.MailerService,
	binanceCryptoParser *parser.BinanceCryptoParser,
	subscriberRepository SubscriberRepository,
	validator *validator.Validate,
) *MailerHandler {
	return &MailerHandler{
		mailerService:             mailerService,
		binanceCryptoParser:       binanceCryptoParser,
		subscriberRepository:      subscriberRepository,
		validator:                 validator,
		exchangeRateBaseCurrency:  model.Currency(os.Getenv("BASE_CURRENCY")),
		exchangeRateQuoteCurrency: model.Currency(os.Getenv("QUOTE_CURRENCY")),
	}
}

func (h *MailerHandler) SendExchangeRate(c *fiber.Ctx) error {
	value, err := h.binanceCryptoParser.GetExchangeRateValue(h.exchangeRateBaseCurrency, h.exchangeRateQuoteCurrency)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	err = h.mailerService.SendValueToAllEmails(strconv.FormatFloat(value, 'f', 2, 64))
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *MailerHandler) Subscribe(c *fiber.Ctx) error {
	var payload subscribeDTO

	if err := c.BodyParser(&payload); err != nil {
		log.Println(err)
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
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

type subscribeDTO struct {
	Email string `validate:"required,email"`
}
