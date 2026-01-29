package handler

import (
	"go-elastic/service"

	"github.com/gofiber/fiber/v2"
)

func HelloHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	msg, err := service.HelloService(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": msg,
	})
}
