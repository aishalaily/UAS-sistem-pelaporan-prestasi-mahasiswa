package config

import (
	"uas-go/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewFiberApp() *fiber.App {
	app := fiber.New()

	app.Use(middleware.Logger())

	return app
}
