package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.uber.org/zap"

	"github.com/xoticdsign/auf-citaty-api/cache"
	"github.com/xoticdsign/auf-citaty-api/database"
	"github.com/xoticdsign/auf-citaty-api/errorhandler"
	"github.com/xoticdsign/auf-citaty-api/logging"
	"github.com/xoticdsign/auf-citaty-api/middleware"
	"github.com/xoticdsign/auf-citaty-api/routes"
)

// TODO:
// ??? LEARN AND USE DOCKER
// ??? LEARN AND USE SWAGGER
// ??? LEARN ABOUT CACHES

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = cache.RunRedis()
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
		ServerHeader:  "auf-citaty-api",
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  errorhandler.ErrorHandler,
		AppName:       "auf-citaty-api",
	})

	middleware.GetMiddleware(api)
	routes.GetRoutes(api)

	logging.Logger.Info(
		"Сервер запущен",
		zap.String("Адрес", "0.0.0.0:8080"),
	)

	err = api.Listen("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
