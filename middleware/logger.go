package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		fmt.Printf("[%s] %s %s - %d (%v)\n",
			time.Now().Format(time.RFC3339),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			time.Since(start),
		)

		return err
	}
}
