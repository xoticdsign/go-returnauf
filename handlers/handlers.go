package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"

	"github.com/xoticdsign/auf-citaty-api/database"
	"github.com/xoticdsign/auf-citaty-api/logging"
)

func ListAll(c *fiber.Ctx) error {
	logging.Logger.Info(
		"Запрос получен",
		zap.String("Метод", c.Method()),
		zap.String("Путь", c.Path()),
	)

	quotes := database.ListAll()

	logging.Logger.Info(
		"Ответ отправлен",
		zap.String("Путь", c.Path()),
		zap.Duration("Время обработки", time.Since(c.Locals("time").(time.Time))),
	)

	return c.JSON(quotes)
}

func RandomQuote(c *fiber.Ctx) error {
	logging.Logger.Info(
		"Запрос получен",
		zap.String("Метод", c.Method()),
		zap.String("Путь", c.Path()),
	)

	quote := database.RandomQuote()

	logging.Logger.Info(
		"Ответ отправлен",
		zap.String("Путь", c.Path()),
		zap.Duration("Время обработки", time.Since(c.Locals("time").(time.Time))),
	)

	return c.JSON(quote)
}

func QuoteID(c *fiber.Ctx) error {
	logging.Logger.Info(
		"Запрос получен",
		zap.String("Метод", c.Method()),
		zap.String("Путь", c.Path()),
	)

	id := c.Params("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	quote, err := database.QuoteID(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	logging.Logger.Info(
		"Ответ отправлен",
		zap.String("Путь", c.Path()),
		zap.Duration("Время обработки", time.Since(c.Locals("time").(time.Time))),
	)

	return c.JSON(quote)
}
