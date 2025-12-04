package service

import (
	"github.com/gofiber/fiber/v2"
	"uas-go/app/model"
	"uas-go/app/repository"
	"uas-go/utils"
)

type AuthService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewAuthService() AuthService {
	return AuthService{}
}

func (s AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "username incorrect"})
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "password incorrect"})
	}

	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "role not found"})
	}

	permissions, _ := s.permissionRepo.GetPermissionsByRole(user.RoleID)

	userResp := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Email:       user.Email,
		Role:        role.Name,
		Permissions: permissions,
	}

	token, err := utils.GenerateToken(userResp)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(model.LoginResponse{
		Token:        token,
		RefreshToken: "",
		User:         userResp,
	})
}


