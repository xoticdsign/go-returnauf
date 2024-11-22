package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/xoticdsign/auf-citaty-api/docs"
	"github.com/xoticdsign/auf-citaty-api/handlers"
)

func GetRoutes(api *fiber.App) {
	api.Get("/swagger/*", swagger.HandlerDefault)
	api.Get("/", handlers.ListAll)
	api.Get("/random", handlers.RandomQuote)
	api.Get("/:id", handlers.QuoteID)
}
