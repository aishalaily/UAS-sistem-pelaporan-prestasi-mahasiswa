package middleware

import "github.com/gofiber/fiber/v2"

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		if role == nil || role.(string) != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "admin only",
			})
		}

		return c.Next()
	}
}

func LecturerOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		if role == nil || role.(string) != "dosen" {
			return c.Status(403).JSON(fiber.Map{
				"error": "lecturer only",
			})
		}

		return c.Next()
	}
}

func StudentOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		if role == nil || role.(string) != "mahasiswa" {
			return c.Status(403).JSON(fiber.Map{
				"error": "student only",
			})
		}

		return c.Next()
	}
}
