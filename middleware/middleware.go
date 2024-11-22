package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"github.com/xoticdsign/auf-citaty-api/errorhandler"
)

func GetMiddleware(api *fiber.App) {
	api.Use(favicon.New(favicon.ConfigDefault))
	api.Use(keyauth.New(keyauth.Config{
		Next:         authFiler,
		ErrorHandler: errorhandler.ErrorHandler,
		KeyLookup:    "query:" + "auf-citaty-key",
		Validator:    keyauthValidator,
	}))
}

func authFiler(c *fiber.Ctx) bool {
	path := c.Path()
	if strings.Contains(path, "swagger") {
		return true
	}
	return false
}

func keyauthValidator(c *fiber.Ctx, key string) (bool, error) {
	timeCounter := time.Now()

	c.Locals("time", timeCounter)

	apiKey := os.Getenv("AUF_CITATY_KEY")

	hRealKey := sha256.Sum256([]byte(apiKey))
	hKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hRealKey[:], hKey[:]) == 1 {
		return true, nil
	}
	return false, fiber.ErrUnauthorized
}
