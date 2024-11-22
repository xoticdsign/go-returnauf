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

// List all quotes
//
//	@description	Возвращает полный список цитат, хранящихся в базе данных. Полезно для получения всех доступных данных для анализа, отображения или других операций. Цитаты возвращаются в формате JSON.
//
//	@id				list-all
//	@tags			Операции с цитатами
//
//	@summary		Предоставляет все цитаты
//	@produce		json
//	@security		KeyAuth
//	@success		200	{object}	database.Quote
//	@failure		401	{object}	errorhandler.Error	Происходит,	если	не			был				предоставлен	ключ	API
//	@failure		500	{object}	errorhandler.Error	Происходит,	если	произошла	неопределенная	ошибка
//	@router			/ [get]
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

// Random quote
//
//	@description	Возвращает случайную цитату из базы данных. Если цитата отсутствует в кэше, то она извлекается из базы данных, добавляется в кэш и возвращается пользователю. Позволяет отображать динамическое содержимое, не перегружая базу данных. Случайность обеспечивается генератором случайных чисел.
//
//	@id				random-quote
//	@tags			Операции с цитатами
//
//	@summary		Предоставляет случайную цитату
//	@produce		json
//	@security		KeyAuth
//	@success		200	{object}	database.Quote
//	@failure		401	{object}	errorhandler.Error	Происходит,	если	не			был				предоставлен	ключ	API
//	@failure		500	{object}	errorhandler.Error	Происходит,	если	произошла	неопределенная	ошибка
//	@router			/random [get]
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

// Quote by ID
//
//	@description	Возвращает цитату по её уникальному идентификатору (ID). Если цитата не найдена в кэше, происходит обращение к базе данных. Полученная цитата затем сохраняется в кэш для ускорения последующих запросов. Если запрошенного ID нет в базе данных, возвращается ошибка.
//
//	@id				quote-id
//	@tags			Операции с цитатами
//
//	@summary		Предоставляет цитату по заданному ID
//	@produce		json
//	@param			id	path	string	false	"Позволяет указать ID цитаты"	example(105)
//	@security		KeyAuth
//	@success		200	{object}	database.Quote
//	@failure		401	{object}	errorhandler.Error	Происходит,	если	не				был				предоставлен	ключ	API
//	@failure		404	{object}	errorhandler.Error	Происходит,	если	запрашиваемой	цитаты			не				существует
//	@failure		500	{object}	errorhandler.Error	Происходит,	если	произошла		неопределенная	ошибка
//	@router			/{id} [get]
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
