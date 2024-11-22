package errhandling

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"go.uber.org/zap"

	"github.com/xoticdsign/auf-citaty/models/responses"
	"github.com/xoticdsign/auf-citaty/utils/logging"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == keyauth.ErrMissingOrMalformedAPIKey {
		logging.Logger.Error(
			fiber.ErrUnauthorized.Message,
			zap.Int("Status", fiber.StatusUnauthorized),
			zap.String("Method", c.Method()),
			zap.String("Path", c.Path()),
		)

		return c.Status(fiber.StatusUnauthorized).JSON(responses.Error{
			Code:    fiber.StatusUnauthorized,
			Message: fiber.ErrUnauthorized.Message,
		})
	}

	var e *fiber.Error

	if errors.As(err, &e) {
		logging.Logger.Error(
			e.Message,
			zap.Int("Status", e.Code),
			zap.String("Method", c.Method()),
			zap.String("Path", c.Path()),
		)

		return c.Status(e.Code).JSON(responses.Error{
			Code:    e.Code,
			Message: e.Message,
		})
	}
	logging.Logger.Error(
		err.Error(),
		zap.Int("Status", fiber.StatusInternalServerError),
		zap.String("Method", c.Method()),
		zap.String("Path", c.Path()),
	)

	return c.Status(fiber.StatusInternalServerError).JSON(responses.Error{
		Code:    fiber.StatusInternalServerError,
		Message: err.Error(),
	})
}
