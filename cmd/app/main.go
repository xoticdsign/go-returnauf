package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"go.uber.org/zap"

	_ "github.com/xoticdsign/auf-citaty/docs"
	"github.com/xoticdsign/auf-citaty/internal/cache"
	"github.com/xoticdsign/auf-citaty/internal/database"
	"github.com/xoticdsign/auf-citaty/internal/handlers"
	"github.com/xoticdsign/auf-citaty/internal/middleware"
	"github.com/xoticdsign/auf-citaty/utils/errhandling"
	"github.com/xoticdsign/auf-citaty/utils/logging"
)

// General description
//
//	@title						Auf Citaty API
//	@version					1.0.0
//	@description				TODO
//	@contact.name				xoti$
//	@contact.url				https://t.me/xoticdsign
//	@contact.email				xoticdollarsign@outlook.com
//	@license.name				MIT
//	@license.url				https://mit-license.org/
//	@host						127.0.0.1:8080
//	@BasePath					/
//	@produce					json
//	@schemes					http
//
//	@securitydefinitions.apikey	KeyAuth
//	@in							query
//	@name						auf-citaty-key
func main() {
	godotenv.Load()

	err := cache.RunRedis()
	if err != nil {
		log.Fatal(err)
	}

	err = logging.RunZap()
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := database.RunGORM()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		ServerHeader:  "auf-citaty",
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  errhandling.ErrorHandler,
		AppName:       "auf-citaty",
	})

	middleware.GetMiddleware(app)

	handler := &handlers.Handlers{DB: gormDB}

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/", handler.ListAll)
	app.Get("/random", handler.RandomQuote)
	app.Get("/:id", handler.QuoteID)

	logging.Logger.Info(
		"Сервер запущен",
		zap.String("Address", os.Getenv("SERVER_ADDRESS")),
	)

	err = app.Listen(os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
}
