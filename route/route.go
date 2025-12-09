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

	auth := api.Group("/auth")
	auth.Post("/login", service.Login)
	auth.Get("/profile", middleware.AuthRequired(), service.GetProfile)

	users := api.Group("/users")
	users.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), service.CreateUser)

	ach := api.Group("/achievements", middleware.AuthRequired())

	ach.Get("/", service.GetAchievements)          
	ach.Get("/:id", service.GetAchievementDetail) 

	ach.Post("/", middleware.MahasiswaOnly(), service.SubmitAchievement)                 
	ach.Put("/:id", middleware.MahasiswaOnly(), service.UpdateAchievement)                
	ach.Post("/:id/submit", middleware.MahasiswaOnly(), service.SubmitForVerification)    
	ach.Delete("/:id", middleware.MahasiswaOnly(), service.DeleteAchievement)            
}
