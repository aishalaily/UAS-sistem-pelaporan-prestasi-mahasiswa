package service

import (
	"github.com/gofiber/fiber/v2"
	"uas-go/app/repository"

	"fmt"
)

type ProfileService struct {
	UserRepo repository.UserRepository
}

func (s ProfileService) GetProfile(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(string)

	profile, err := s.UserRepo.FindByID(userID)
	if err != nil {
		fmt.Println("DEBUG FIND PROFILE ERROR:", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to load profile"})
	}

	return c.JSON(profile)
}
