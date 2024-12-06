package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "github.com/xoticdsign/auf-citaty/docs"
	"github.com/xoticdsign/auf-citaty/internal/cache"
	"github.com/xoticdsign/auf-citaty/internal/database"
	"github.com/xoticdsign/auf-citaty/internal/handlers"
	"github.com/xoticdsign/auf-citaty/internal/logging"
	"github.com/xoticdsign/auf-citaty/internal/middleware"
)

// Общее описание
//
// @title                      Auf Citaty API
// @version                    1.0.0
// @description                TODO
// @contact.name               xoti$
// @contact.url                https://t.me/xoticdsign
// @contact.email              xoticdollarsign@outlook.com
// @license.name               MIT
// @license.url                https://mit-license.org/
// @host                       127.0.0.1:8080
// @BasePath                   /
// @produce                    json
// @schemes                    http
//
// @securitydefinitions.apikey KeyAuth
// @in                         query
// @name                       auf-citaty-key
func main() {
	godotenv.Load()

	serverAddr := os.Getenv("SERVER_ADDRESS")
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	dbAddr := os.Getenv("DB_ADDRESS")

	cache, err := cache.RunRedis(redisAddr, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.RunZap()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.RunGORM(dbAddr)
	if err != nil {
		log.Fatal(err)
	}

	dependencies := &handlers.Dependencies{
		DB:     db,
		Cache:  cache,
		Logger: logger,
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

	err = app.Listen(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
}
