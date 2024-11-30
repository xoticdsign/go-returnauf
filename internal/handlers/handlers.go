package handlers

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"github.com/xoticdsign/auf-citaty/internal/cache"
	"github.com/xoticdsign/auf-citaty/internal/database"
	"github.com/xoticdsign/auf-citaty/internal/logging"
	"github.com/xoticdsign/auf-citaty/models/responses"
)

type Dependencies struct {
	DB     database.Queuer
	Cache  cache.Cacher
	Logger logging.Logger
}

func (d *Dependencies) Error(c *fiber.Ctx, err error) error {
	if err == keyauth.ErrMissingOrMalformedAPIKey {
		d.Logger.Error(fiber.ErrUnauthorized.Message, c)

		return c.Status(fiber.StatusUnauthorized).JSON(responses.Error{
			Code:    fiber.StatusUnauthorized,
			Message: fiber.ErrUnauthorized.Message,
		})
	}

	var e *fiber.Error

	if errors.As(err, &e) {
		d.Logger.Error(e.Message, c)

		return c.Status(e.Code).JSON(responses.Error{
			Code:    e.Code,
			Message: e.Message,
		})
	}
	d.Logger.Error(err.Error(), c)

	return c.Status(fiber.StatusInternalServerError).JSON(responses.Error{
		Code:    fiber.StatusInternalServerError,
		Message: err.Error(),
	})
}

// List all quotes
//
// @description Возвращает полный список цитат, хранящихся в базе данных. Полезно для получения всех доступных данных для анализа, отображения или других операций. Цитаты возвращаются в формате JSON.
//
// @id          list-all
// @tags        Операции с цитатами
//
// @summary     Предоставляет все цитаты
// @produce     json
// @security    KeyAuth
// @success     200 {object} responses.Quote Стандартный ответ
// @failure     401 {object} responses.Error Происходит, если не        был            предоставлен ключ API
// @failure     500 {object} responses.Error Происходит, если произошла неопределенная ошибка
// @router      / [get]
func (d *Dependencies) ListAll(c *fiber.Ctx) error {
	d.Logger.Info("Обращение к базе данных", c)

	quotes := d.DB.ListAll()

	d.Logger.Info("Обработан запрос", c)

	return c.JSON(quotes)
}

// Random quote
//
// @description Возвращает случайную цитату из базы данных. Если цитата отсутствует в кэше, то она извлекается из базы данных, добавляется в кэш и возвращается пользователю. Позволяет отображать динамическое содержимое, не перегружая базу данных. Случайность обеспечивается генератором случайных чисел.
//
// @id          random-quote
// @tags        Операции с цитатами
//
// @summary     Предоставляет случайную цитату
// @produce     json
// @security    KeyAuth
// @success     200 {object} responses.Quote Стандартный ответ
// @failure     401 {object} responses.Error Происходит, если не        был            предоставлен ключ API
// @failure     500 {object} responses.Error Происходит, если произошла неопределенная ошибка
// @router      /random [get]
func (d *Dependencies) RandomQuote(c *fiber.Ctx) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randInt := rand.Intn(201)

	id := strconv.Itoa(randInt)

	quote, errStr := d.Cache.Get(id)
	if errStr == "Failed" {
		d.Logger.Error("Не удалось достать кэш", c)
	}
	if errStr == "Nil" {
		d.Logger.Warn("Кэш отсутствует", c)
		d.Logger.Info("Обращение к базе данных", c)

		quote, _ := d.DB.GetQuote(id)

		err := d.Cache.Set(id, quote.Quote, time.Minute*1)
		if err != nil {
			d.Logger.Error("Не удалось кэшировать данные", c)
		}

		d.Logger.Info("Данные добавлены в кэш", c)
		d.Logger.Info("Обработан запрос", c)

		return c.JSON(quote)
	}

	d.Logger.Info("Данные получены из кэша", c)
	d.Logger.Info("Обработан запрос", c)

	return c.JSON(responses.Quote{
		ID:    randInt,
		Quote: quote,
	})
}

// Quote by ID
//
// @description Возвращает цитату по её уникальному идентификатору (ID). Если цитата не найдена в кэше, происходит обращение к базе данных. Полученная цитата затем сохраняется в кэш для ускорения последующих запросов. Если запрошенного ID нет в базе данных, возвращается ошибка.
//
// @id          quote-id
// @tags        Операции с цитатами
//
// @summary     Предоставляет цитату по заданному ID
// @produce     json
// @param       id path string false "Позволяет указать ID цитаты" example(105)
// @security    KeyAuth
// @success     200 {object} responses.Quote Стандартный ответ
// @failure     401 {object} responses.Error Происходит, если не            был            предоставлен ключ API
// @failure     404 {object} responses.Error Происходит, если запрашиваемой цитаты         не           существует
// @failure     500 {object} responses.Error Происходит, если произошла     неопределенная ошибка
// @router      /{id} [get]
func (d *Dependencies) QuoteID(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	quote, errStr := d.Cache.Get(id)
	if errStr == "Failed" {
		d.Logger.Error("Не удалось достать кэш", c)
	}
	if errStr == "Nil" {
		d.Logger.Warn("Кэш отсутствует", c)
		d.Logger.Info("Обращение к базе данных", c)

		quote, err := d.DB.GetQuote(id)
		if err != nil {
			return fiber.ErrNotFound
		}

		err = d.Cache.Set(id, quote.Quote, time.Minute*1)
		if err != nil {
			d.Logger.Error("Не удалось кэшировать данные", c)
		}

		d.Logger.Info("Данные добавлены в кэш", c)
		d.Logger.Info("Обработан запрос", c)

		return c.JSON(quote)
	}

	d.Logger.Info("Данные получены из кэша", c)
	d.Logger.Info("Обработан запрос", c)

	return c.JSON(responses.Quote{
		ID:    idInt,
		Quote: quote,
	})
}
