package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/google/uuid"

	"github.com/xoticdsign/auf-citaty/config"
	"github.com/xoticdsign/auf-citaty/internal/cache"
	"github.com/xoticdsign/auf-citaty/internal/database"
	"github.com/xoticdsign/auf-citaty/internal/handlers"
	"github.com/xoticdsign/auf-citaty/internal/logging"
	"github.com/xoticdsign/auf-citaty/internal/middleware"
	"github.com/xoticdsign/auf-citaty/internal/utils"
)

// Инициализирует приложение
func InitApp(conf config.Config) (*fiber.App, error) {
	cache, err := cache.RunRedis(conf.RedisAddr, conf.RedisPassword)
	if err != nil {
		return nil, err
	}

	logger, err := logging.RunZap()
	if err != nil {
		return nil, err
	}

	db, err := database.RunGORM(conf.DBAddr)
	if err != nil {
		return nil, err
	}

	dependencies := &handlers.Dependencies{
		DB:      db,
		Cache:   cache,
		Logger:  logger,
		Support: &utils.Support{},
	}

	app := fiber.New(fiber.Config{
		ServerHeader:  "auf-citaty",
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  dependencies.Error,
		AppName:       "auf-citaty",
	})

	app.Use(favicon.New(favicon.ConfigDefault))
	app.Use(requestid.New(requestid.Config{
		Generator:  uuid.NewString,
		ContextKey: "uuid",
	}))
	app.Use(keyauth.New(keyauth.Config{
		Next:         middleware.AuthFiler,
		ErrorHandler: dependencies.Error,
		KeyLookup:    "query:auf-citaty-key",
		Validator:    middleware.KeyauthValidator,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/", dependencies.ListAll)
	app.Get("/random", dependencies.RandomQuote)
	app.Get("/:id", dependencies.QuoteID)

	return app, nil
}
