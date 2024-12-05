package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Фильтрует маршруты для аутентификации
func AuthFiler(c *fiber.Ctx) bool {
	path := c.Path()
	if strings.Contains(path, "swagger") {
		return true
	}
	return false
}

// Проверяет ключ API
func KeyauthValidator(c *fiber.Ctx, key string) (bool, error) {
	apiKey := os.Getenv("AUF_CITATY_KEY")

	hRealKey := sha256.Sum256([]byte(apiKey))
	hKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hRealKey[:], hKey[:]) == 1 {
		return true, nil
	}
	return false, fiber.ErrUnauthorized
}
