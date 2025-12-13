package middleware

import (
	"strings"
	"uas-go/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"error": "authorization header missing or invalid",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		role := strings.ToLower(strings.TrimSpace(claims.RoleName))
		role = strings.ReplaceAll(role, " ", "_")
		c.Locals("role", role)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}
