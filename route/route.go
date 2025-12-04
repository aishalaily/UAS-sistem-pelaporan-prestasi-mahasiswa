package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-go/app/service"
	"uas-go/middleware"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// sementara kosong, hanya test
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	authService := service.NewAuthService()
	profileService := service.ProfileService{}

	api.Post("/login", func(c *fiber.Ctx) error { return authService.Login(c) })
	api.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
    return profileService.GetProfile(c)
})
	// protected := api.Group("", middleware.AuthRequired())
	// admin := protected.Group("/admin", middleware.AdminOnly())

}