package service

import (
    "database/sql"
    "os"

    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"

    "uas-go/app/repository"
)

func LoginService(c *fiber.Ctx) error {
    type LoginRequest struct {
        Username    string `json:"username"`
        Password string `json:"password"`
    }

    req := new(LoginRequest)
    if err := c.BodyParser(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Invalid JSON",
        })
    }

    user, err := repository.FindUserByUsername(req.Username)
    if err == sql.ErrNoRows {
        return c.Status(401).JSON(fiber.Map{
            "success": false,
            "message": "Username tidak ditemukan",
        })
    }
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Server error",
        })  
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
        return c.Status(401).JSON(fiber.Map{
            "success": false,
            "message": "Password salah",
        })
    }

    // generate JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role_id": user.RoleID,
    })

    jwtSecret := []byte(os.Getenv("JWT_SECRET"))
    tokenString, _ := token.SignedString(jwtSecret)

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Login berhasil",
        "token":   tokenString,
    })
}

func ProfileService(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	username := c.Locals("username").(string)
	role := c.Locals("name").(string)

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "User tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":        user.ID,
			"username":  username,
			"name":      role,
		},
	})
}
