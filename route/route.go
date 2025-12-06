package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-go/app/service"
	"uas-go/middleware"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// hanya test
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	api.Post("/auth/login", service.Login)
	api.Get("/auth/profile", middleware.AuthRequired(), service.GetProfile)

	api.Post("/users", middleware.AuthRequired(), middleware.AdminOnly(), service.CreateUser)

	api.Post("/achievement", middleware.AuthRequired(), middleware.MahasiswaOnly(), service.SubmitAchievement)
	api.Delete("/achievement/:id", middleware.AuthRequired(), middleware.MahasiswaOnly(), service.DeleteAchievement)
	api.Post("/achievement/:id", middleware.AuthRequired(), middleware.MahasiswaOnly(), service.SubmitForVerification)
}
