package handler

import (
	"github.com/gofiber/fiber/v2"
)

type ResponsePresenter interface {
	PresentExchangeRate(c *fiber.Ctx, rate float64) error
	PresentError(c *fiber.Ctx, statusCode int, err error) error
}

type JSONPresenter struct{}

func NewJSONPresenter() ResponsePresenter {
	return &JSONPresenter{}
}

func (p JSONPresenter) PresentExchangeRate(c *fiber.Ctx, rate float64) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"rate": rate})
}

func (p JSONPresenter) PresentError(c *fiber.Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(&fiber.Map{"message": err.Error()})
}
