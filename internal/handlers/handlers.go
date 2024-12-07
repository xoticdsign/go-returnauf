package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"github.com/xoticdsign/auf-citaty/internal/cache"
	"github.com/xoticdsign/auf-citaty/internal/database"
	"github.com/xoticdsign/auf-citaty/internal/logging"
	"github.com/xoticdsign/auf-citaty/internal/utils"
	"github.com/xoticdsign/auf-citaty/models/responses"
)

// Структура, содержащая интерфейсы для инъекции
type Dependencies struct {
	DB      database.Queuer
	Cache   cache.Cacher
	Logger  logging.Logger
	Support utils.Supporter
}

// Получает контекст и ошибку, а затем форматирует все в JSON
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
		er, ok := responses.ErrDictionary[e.Code]
		if !ok {
			d.Logger.Warn("Необработанная ошибка: "+er.Message, c)

			return c.Status(fiber.StatusInternalServerError).JSON(responses.Error{
				Code:    fiber.StatusInternalServerError,
				Message: fiber.ErrInternalServerError.Message,
			})
		}
		d.Logger.Error(er.Message, c)

		return c.Status(er.Code).JSON(responses.Error{
			Code:    er.Code,
			Message: er.Message,
		})
	}
	d.Logger.Error(fiber.ErrInternalServerError.Message, c)

	return c.Status(fiber.StatusInternalServerError).JSON(responses.Error{
		Code:    fiber.StatusInternalServerError,
		Message: fiber.ErrInternalServerError.Message,
	})
}

// @description Возвращает полный список цитат, хранящихся в базе данных. Полезно для получения всех доступных данных для анализа, отображения или других операций. Цитаты возвращаются в формате JSON.
//
// @id          list-all
// @tags        Операции с цитатами
//
// @summary     Предоставляет все цитаты
// @produce     json
// @security    KeyAuth
// @success     200 {object} responses.Quote
// @failure     401 {object} responses.Error
// @failure     404 {object} responses.Error
// @failure     405 {object} responses.Error
// @failure     500 {object} responses.Error
// @router      / [get]
func (d *Dependencies) ListAll(c *fiber.Ctx) error {
	quotes, err := d.DB.ListAll()
	if err != nil {
		return fiber.ErrNotFound
	}
	d.Logger.Info("Обработан запрос", c)

	return c.JSON(quotes)
}

// @description Возвращает случайную цитату из базы данных. Если цитата отсутствует в кэше, то она извлекается из базы данных, добавляется в кэш и возвращается пользователю. Позволяет отображать динамическое содержимое, не перегружая базу данных. Случайность обеспечивается генератором случайных чисел.
//
// @id          random-quote
// @tags        Операции с цитатами
//
// @summary     Предоставляет случайную цитату
// @produce     json
// @security    KeyAuth
// @success     200 {object} responses.Quote
// @failure     401 {object} responses.Error
// @failure     404 {object} responses.Error
// @failure     405 {object} responses.Error
// @failure     500 {object} responses.Error
// @router      /random [get]
func (d *Dependencies) RandomQuote(c *fiber.Ctx) error {
	count, err := d.DB.QuotesCount()
	if err != nil {
		return fiber.ErrNotFound
	}

	idInt, id := d.Support.RandInt(count)

	quote, err := d.Cache.Get(id)
	if err != nil {
		quote, err := d.DB.GetQuote(id)
		if err != nil {
			return fiber.ErrNotFound
		}

		err = d.Cache.Set(id, quote.Quote, time.Minute*1)
		if err != nil {
			d.Logger.Error("Не удалось кэшировать данные", c)
		}
		d.Logger.Info("Обработан запрос", c)

		return c.JSON(quote)
	}
	d.Logger.Info("Обработан запрос", c)

	return c.JSON(responses.Quote{
		ID:    idInt,
		Quote: quote,
	})
}

// @description Возвращает цитату по её уникальному идентификатору (ID). Если цитата не найдена в кэше, происходит обращение к базе данных. Полученная цитата затем сохраняется в кэш для ускорения последующих запросов. Если запрошенного ID нет в базе данных, возвращается ошибка.
//
// @id          quote-id
// @tags        Операции с цитатами
//
// @summary     Предоставляет цитату по заданному ID
// @produce     json
// @param       id path string false "Позволяет указать ID цитаты" example(105)
// @security    KeyAuth
// @success     200 {object} responses.Quote
// @failure     401 {object} responses.Error
// @failure     404 {object} responses.Error
// @failure     405 {object} responses.Error
// @failure     500 {object} responses.Error
// @router      /{id} [get]
func (d *Dependencies) QuoteID(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	quote, err := d.Cache.Get(id)
	if err != nil {
		quote, err := d.DB.GetQuote(id)
		if err != nil {
			return fiber.ErrNotFound
		}

		err = d.Cache.Set(id, quote.Quote, time.Minute*1)
		if err != nil {
			d.Logger.Error("Не удалось кэшировать данные", c)
		}
		d.Logger.Info("Обработан запрос", c)

		return c.JSON(quote)
	}
	d.Logger.Info("Обработан запрос", c)

	return c.JSON(responses.Quote{
		ID:    idInt,
		Quote: quote,
	})
}
