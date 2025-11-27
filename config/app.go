package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"uas-go/database"
	"uas-go/route"
)

func SetupApp() *fiber.App {
	LoadEnv()

	database.ConnectPostgres()
	database.ConnectMongo()

	app := fiber.New()

	route.RegisterRoutes(app)

	log.Println("App initialized")
	return app
}
