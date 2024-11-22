package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.uber.org/zap"

	"github.com/xoticdsign/auf-citaty-api/cache"
	"github.com/xoticdsign/auf-citaty-api/database"
	_ "github.com/xoticdsign/auf-citaty-api/docs"
	"github.com/xoticdsign/auf-citaty-api/errorhandler"
	"github.com/xoticdsign/auf-citaty-api/logging"
	"github.com/xoticdsign/auf-citaty-api/middleware"
	"github.com/xoticdsign/auf-citaty-api/routes"
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

	err = database.RunGORM()
	if err != nil {
		log.Fatal(err)
	}

	api := fiber.New(fiber.Config{
		ServerHeader:  "auf-citaty",
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  errorhandler.ErrorHandler,
		AppName:       "auf-citaty",
	})

	middleware.GetMiddleware(api)
	routes.GetRoutes(api)

	addr := os.Getenv("SERVER_ADDRESS")

	logging.Logger.Info(
		"Сервер запущен",
		zap.String("Адрес", addr),
	)

	err = api.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}
}
