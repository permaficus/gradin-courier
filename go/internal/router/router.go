package router

import (
	"courier-technical-test/go/internal/courier"
	"courier-technical-test/go/internal/response"
	"github.com/gofiber/fiber/v2"
)

func New(handler *courier.Handler) *fiber.App {
	app := fiber.New(fiber.Config{AppName: "Courier Technical Test Go"})
	app.Use(func(c *fiber.Ctx) error {
		response.RequestID(c)
		return c.Next()
	})

	api := app.Group("/api")
	couriers := api.Group("/couriers")
	couriers.Get("/", handler.Index)
	couriers.Post("/", handler.Store)
	couriers.Get("/:id", handler.Show)
	couriers.Put("/:id", handler.Update)
	couriers.Delete("/:id", handler.Destroy)
	return app
}
