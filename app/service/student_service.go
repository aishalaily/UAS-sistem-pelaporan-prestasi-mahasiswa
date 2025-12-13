package service

import (
	"uas-go/app/repository"
	"uas-go/database"

	"github.com/gofiber/fiber/v2"
)

func GetStudents(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	data, err := repository.GetAllStudents()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load students",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func GetStudentByID(c *fiber.Ctx) error {
	studentID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	data, err := repository.GetStudentByID(studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	switch role {
	case "admin":

	case "mahasiswa":
		sid, _ := repository.GetStudentIDByUserID(database.PgPool, userID)
		if sid != studentID {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied",
			})
		}

	case "dosen_wali":
		ok, _ := repository.IsStudentUnderAdvisor(
			database.PgPool,
			userID,
			studentID,
		)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "not your advisee",
			})
		}

	default:
		return c.Status(403).JSON(fiber.Map{
			"error": "role not allowed",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	switch role {
	case "admin":

	case "mahasiswa":
		sid, _ := repository.GetStudentIDByUserID(database.PgPool, userID)
		if sid != studentID {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied",
			})
		}

	case "dosen_wali":
		ok, _ := repository.IsStudentUnderAdvisor(
			database.PgPool,
			userID,
			studentID,
		)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "not your advisee",
			})
		}

	default:
		return c.Status(403).JSON(fiber.Map{
			"error": "role not allowed",
		})
	}

	refs, err := repository.GetAchievementsByStudent(database.PgPool, studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load achievements",
		})
	}

	var result []fiber.Map
	for _, ref := range refs {
		mongoData, err := repository.GetAchievementMongo(ref.MongoAchievementID)
		if err != nil {
			continue 
		}

		result = append(result, fiber.Map{
			"reference": fiber.Map{
				"id":         ref.ID,
				"status":     ref.Status,
				"created_at": ref.CreatedAt,
			},
			"details": mongoData,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func UpdateStudentAdvisor(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	studentID := c.Params("id")

	var body struct {
		AdvisorID string `json:"advisor_id"`
	}

	if err := c.BodyParser(&body); err != nil || body.AdvisorID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "advisor_id is required",
		})
	}

	_, err := repository.GetStudentByID(studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	if err := repository.UpdateStudentAdvisor(studentID, body.AdvisorID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update advisor",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "advisor updated",
	})
}

