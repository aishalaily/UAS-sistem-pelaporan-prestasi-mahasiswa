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
		auth.Post("/refresh", service.RefreshToken)
		auth.Post("/logout", middleware.AuthRequired(), service.Logout)

	users := api.Group("/users")
		users.Get("/", middleware.AuthRequired(), middleware.AdminOnly(), service.GetUsers)
		users.Get("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.GetUserDetail)
		users.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), service.CreateUser)
		users.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.UpdateUser)
		users.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.DeleteUser)
		users.Put("/:id/role", middleware.AuthRequired(), middleware.AdminOnly(), service.UpdateUserRole)

	ach := api.Group("/achievements", middleware.AuthRequired())
		ach.Get("/", middleware.RequirePermission("achievement.read"), service.GetAchievements)          
		ach.Get("/:id", middleware.RequirePermission("achievement.read"), service.GetAchievementDetail) 

		ach.Post("/", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.create"), service.SubmitAchievement)                 
		ach.Put("/:id", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.update"), service.UpdateAchievement)                
		ach.Post("/:id/submit", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.update"), service.SubmitForVerification)    
		ach.Delete("/:id", middleware.MahasiswaOnly(), middleware.RequirePermission("achievement.delete"), service.DeleteAchievement) 
		
		ach.Post("/:id/verify", middleware.DosenWaliOnly(), middleware.RequirePermission("achievement.verify"), service.VerifyAchievement)
		ach.Post("/:id/reject", middleware.DosenWaliOnly(), middleware.RequirePermission("achievement.reject"), service.RejectAchievement)

		ach.Get("/:id/history", service.GetAchievementHistory)
		ach.Post("/:id/attachments", middleware.MahasiswaOnly(), service.UploadAchievementAttachment)

	students := api.Group("/students", middleware.AuthRequired())
		students.Get("/", middleware.AdminOnly(), service.GetStudents)
		students.Get("/:id", service.GetStudentByID)
		students.Get("/:id/achievements", service.GetStudentAchievements)
		students.Put("/:id/advisor", middleware.AdminOnly(), service.UpdateStudentAdvisor)
	
	lecturers := api.Group("/lecturers", middleware.AuthRequired(), middleware.AdminOnly())
		lecturers.Get("/", middleware.AdminOnly(), service.GetLecturers)
		lecturers.Get("/:id/advisees", service.GetLecturerAdvisees)

	reports := api.Group("/reports", middleware.AuthRequired())
		reports.Get("/statistics", service.GetAchievementStatistics)
		reports.Get("/student/id", service.GetStudentReport)

}
