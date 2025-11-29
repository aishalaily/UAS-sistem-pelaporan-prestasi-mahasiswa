package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"uas-go/database"
	"uas-go/middleware"
	"uas-go/route"
)

func SetupApp() *fiber.App {
	LoadEnv()

	database.ConnectPostgres()
	database.ConnectMongo()

	app := fiber.New()

	app.Use(middleware.LoggerMiddleware)

	route.RegisterRoutes(app)

	log.Println("App initialized")

	return app
}
