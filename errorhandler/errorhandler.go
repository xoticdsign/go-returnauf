package errorhandler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"go.uber.org/zap"

	"github.com/xoticdsign/auf-citaty-api/logging"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	logging.Logger.Error(
		"Ошибка",
		zap.String("Метод", c.Method()),
		zap.String("Путь", c.Path()),
		zap.String("Детали ошибки", err.Error()),
	)

	if err == keyauth.ErrMissingOrMalformedAPIKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: fiber.ErrUnauthorized.Message,
		})
	}

	var e *fiber.Error

	if errors.As(err, &e) {
		return c.Status(e.Code).JSON(e)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
		Code:    fiber.StatusInternalServerError,
		Message: err.Error(),
	})
}
