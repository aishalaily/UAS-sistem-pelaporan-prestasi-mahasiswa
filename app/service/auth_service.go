package service

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"uas-go/app/model"
	"uas-go/app/repository"
	"uas-go/utils"
)

func Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Request tidak valid",
		})
	}

	user, passwordHash, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "username is incorrect",
		})
	}


	if !utils.CheckPassword(passwordHash, req.Password) {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "password is incorrect",
		})
	}

	roleName, err := repository.GetRoleName(user.RoleID)
	if err != nil {
		roleName = "unknown"
	}

	roleReadable := strings.Title(strings.ReplaceAll(roleName, "_", " "))

	permissions, _ := repository.GetPermissionsByRole(user.RoleID)

	userResp := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        roleReadable,
		Permissions: permissions,
	}

	token, err := utils.GenerateToken(userResp)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed generate token",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        token,
			"refreshToken": "",
			"user":         userResp,
		},
	})
}


func GetProfile(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	user, err := repository.GetUserByID(userID.(string))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	roleName, err := repository.GetRoleName(user.RoleID)
	if err != nil {
		roleName = "unknown"
	}

	roleReadable := strings.Title(strings.ReplaceAll(roleName, "_", " "))

	permissions, _ := repository.GetPermissionsByRole(user.RoleID)
	if permissions == nil {
		permissions = []string{}
	}

	resp := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Role:        roleReadable,
		Permissions: permissions,
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user": resp,
		},
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var body struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&body); err != nil || body.Token == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "token is required",
		})
	}

	claims, err := utils.ParseToken(body.Token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	userResp := model.UserResponse{
		ID:          claims.UserID,
		Username:    claims.Username,
		Role:        claims.RoleName,
		Permissions: claims.Permissions,
	}

	newToken, err := utils.GenerateToken(userResp)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": newToken,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "logged out",
	})
}
