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

	ach.Get("/", middleware.RequirePermission("achievement.read"), service.GetAchievements)          
	ach.Get("/:id", middleware.RequirePermission("achievement.read"), service.GetAchievementDetail) 

	ach.Post("/", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.create"), service.SubmitAchievement)                 
	ach.Put("/:id", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.update"), service.UpdateAchievement)                
	ach.Post("/:id/submit", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.update"), service.SubmitForVerification)    
	ach.Delete("/:id", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.delete"), service.DeleteAchievement) 
	
	ach.Post("/:id/verify", middleware.DosenWaliOnly(), middleware.RequirePermission("achievement.verify"), service.VerifyAchievement)
	ach.Post("/:id/reject", middleware.DosenWaliOnly(), middleware.RequirePermission("achievement.reject"), service.RejectAchievement)
}
