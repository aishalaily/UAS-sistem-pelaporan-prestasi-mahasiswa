package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-go/app/service"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// sementara kosong, hanya test
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	api.Post("/login", func(c *fiber.Ctx) error { return service.LoginService(c) })
	api.Get("/profile", func(c *fiber.Ctx) error { return service.ProfileService(c) })
}