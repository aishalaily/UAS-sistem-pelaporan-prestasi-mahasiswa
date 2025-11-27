package route

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// sementara kosong, hanya test
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})
}
