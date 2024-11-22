package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/xoticdsign/auf-citaty/internal/handlers"
)

func GetRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/", handlers.ListAll)
	app.Get("/random", handlers.RandomQuote)
	app.Get("/:id", handlers.QuoteID)
}
