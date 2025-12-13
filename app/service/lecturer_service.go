package service

import (
	"uas-go/app/repository"

	"github.com/gofiber/fiber/v2"
)

func GetLecturers(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	data, err := repository.GetAllLecturers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load lecturers",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	switch role {
	case "admin":

	case "dosen_wali":
		lect, err := repository.GetLecturerByUserID(userID)
		if err != nil || lect.ID != lecturerID {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied",
			})
		}


	default:
		return c.Status(403).JSON(fiber.Map{
			"error": "role not allowed",
		})
	}

	data, err := repository.GetStudentsByAdvisor(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load advisees",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}
