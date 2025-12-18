package service

import (
	"uas-go/database"
	"uas-go/app/repository"

	"github.com/gofiber/fiber/v2"
)

// GetLecturers godoc
// @Summary Get all lecturers
// @Description Admin can view all lecturers
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lecturers [get]
func GetLecturers(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	data, err := repository.GetAllLecturers(database.PgPool)
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

// GetLecturerAdvisees godoc
// @Summary Get lecturer advisees
// @Description Get students supervised by lecturer (Admin or Lecturer himself)
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lecturers/{id}/advisees [get]
func GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	if role == "dosen_wali" {
		lect, err := repository.GetLecturerByUserID(userID)
		if err != nil || lect.ID != lecturerID {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied",
			})
		}
	} else if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "role not allowed",
		})
	}

	data, err := repository.GetStudentsByAdvisor(
		database.PgPool,
		lecturerID,
	)
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
