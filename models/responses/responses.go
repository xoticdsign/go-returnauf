package responses

import (
	"github.com/gofiber/fiber/v2"
)

// Структура для возврата цитаты
type Quote struct {
	ID    int    `gorm:"type:BIGINT NOT NULL PRIMARY KEY"`
	Quote string `gorm:"type:VARCHAR NOT NULL"`
}

// Структура для возврата ошибки
type Error struct {
	Code    int
	Message string
}

// Словарь ошибок
var ErrDictionary = map[int]Error{
	401: {
		Code:    fiber.StatusUnauthorized,
		Message: fiber.ErrUnauthorized.Message,
	},
	404: {
		Code:    fiber.StatusNotFound,
		Message: fiber.ErrNotFound.Message,
	},
	405: {
		Code:    fiber.StatusMethodNotAllowed,
		Message: fiber.ErrMethodNotAllowed.Message,
	},
	500: {
		Code:    fiber.StatusInternalServerError,
		Message: fiber.ErrInternalServerError.Message,
	},
}
