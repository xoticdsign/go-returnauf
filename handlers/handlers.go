package handlers

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"

	"github.com/xoticdsign/auf-citaty-api/cache"
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

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randInt := rand.Intn(201)

	id := strconv.Itoa(randInt)

	quote, err := cache.Cache.Get(context.Background(), id).Result()
	if err == redis.Nil {
		quote, _ := database.GetQoute(id)

		cache.Cache.Set(context.Background(), id, quote.Quote, time.Minute*1)

		logging.Logger.Info(
			"Ответ отправлен",
			zap.String("Путь", c.Path()),
			zap.Duration("Время обработки", time.Since(c.Locals("time").(time.Time))),
		)

		return c.JSON(quote)
	}

	logging.Logger.Info(
		"Ответ отправлен",
		zap.String("Путь", c.Path()),
		zap.Duration("Время обработки", time.Since(c.Locals("time").(time.Time))),
	)

	return c.JSON(fiber.Map{
		"ID":    id,
		"Quote": quote,
	})
}

func QuoteID(c *fiber.Ctx) error {
	logging.Logger.Info(
		"Запроc получен",
		zap.String("Метод", c.Method()),
		zap.String("Путь", c.Path()),
	)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	quote, err := cache.Cache.Get(context.Background(), id).Result()
	if err == redis.Nil {
		quote, err := database.GetQoute(id)
		if err != nil {
			return fiber.ErrNotFound
		}

		cache.Cache.Set(context.Background(), id, quote.Quote, time.Minute*1)

		logging.Logger.Info(
			"Ответ отправлен",
			zap.String("Путь", c.Path()),
			zap.Duration("Время обработки", time.Since(c.Locals("time").(time.Time))),
		)

		return c.JSON(quote)
	}

	logging.Logger.Info(
		"Ответ отправлен",
		zap.String("Путь", c.Path()),
		zap.Duration("Время обработки", time.Since(c.Locals("time").(time.Time))),
	)

	return c.JSON(fiber.Map{
		"ID":    idInt,
		"Quote": quote,
	})
}
