package middleware

import "github.com/gofiber/fiber/v2"

func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")

		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthenticated",
			})
		}

		roleStr := userRole.(string)

		for _, allowed := range roles {
			if roleStr == allowed {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "You do not have permission",
		})
	}
}
