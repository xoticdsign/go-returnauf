package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/google/uuid"

	"github.com/xoticdsign/returnauf/config"
	"github.com/xoticdsign/returnauf/internal/cache"
	"github.com/xoticdsign/returnauf/internal/database"
	"github.com/xoticdsign/returnauf/internal/handlers"
	"github.com/xoticdsign/returnauf/internal/logging"
	"github.com/xoticdsign/returnauf/internal/middleware"
	"github.com/xoticdsign/returnauf/internal/utils"
)

// Инициализирует приложение
func InitApp(conf config.Config) (*fiber.App, error) {
	Cache, err := cache.RunRedis(conf.RedisAddr, conf.RedisPassword)
	if err != nil {
		return nil, err
	}

	Log, err := logging.RunZap()
	if err != nil {
		return nil, err
	}

	DB, err := database.RunGORM(conf.DBAddr)
	if err != nil {
		return nil, err
	}

	dependencies := &handlers.Dependencies{
		DB:      DB,
		Cache:   Cache,
		Logger:  Log,
		Support: &utils.Support{},
	}

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  dependencies.Error,
		AppName:       "returnauf",
	})

	app.Use(favicon.New(favicon.ConfigDefault))
	app.Use(requestid.New(requestid.Config{
		Generator:  uuid.NewString,
		ContextKey: "uuid",
	}))
	app.Use(keyauth.New(keyauth.Config{
		Next:         middleware.AuthFiler,
		ErrorHandler: dependencies.Error,
		KeyLookup:    "query:returnauf-key",
		Validator:    middleware.KeyauthValidator,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/", dependencies.ListAll)
	app.Get("/random", dependencies.RandomQuote)
	app.Get("/:id", dependencies.QuoteID)

	return app, nil
}
