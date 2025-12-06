package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied: role not found",
			})
		}

		if role.(string) != requiredRole {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied: insufficient role",
			})
		}

		return c.Next()
	}
}

func RequirePermission(required string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		perms := c.Locals("permissions")
		if perms == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied: no permissions found",
			})
		}

		permissionList := perms.([]string)

		for _, p := range permissionList {
			if p == required {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "access denied: missing permission",
		})
	}
}

func AdminOnly() fiber.Handler {
	return RequireRole("admin")
}

func DosenWaliOnly() fiber.Handler {
	return RequireRole("dosen_wali")
}

func MahasiswaOnly() fiber.Handler {
	return RequireRole("mahasiswa")
}
