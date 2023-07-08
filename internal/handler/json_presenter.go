package handler

import (
	"github.com/gofiber/fiber/v2"
)

type ResponsePresenter interface {
	PresentExchangeRate(c *fiber.Ctx, rate float64) error
	PresentError(c *fiber.Ctx, err error) error
}

type JSONPresenter struct{}

func NewJSONPresenter() ResponsePresenter {
	return &JSONPresenter{}
}

func (p JSONPresenter) PresentExchangeRate(c *fiber.Ctx, rate float64) error {
	return c.JSON(rate)
}

func (p JSONPresenter) PresentError(c *fiber.Ctx, err error) error {
	return c.JSON(&fiber.Map{"message": err.Error()})
}
