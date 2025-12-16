package service

import (
	"uas-go/app/repository"
	"uas-go/database"
	"uas-go/app/model"

	"github.com/gofiber/fiber/v2"
)

func GetStudents(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	data, err := repository.GetAllStudents(database.PgPool)
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
	studentPK := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	student, err := repository.GetStudentByID(database.PgPool, studentPK)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	switch role {
	case "admin":
		// full access

	case "mahasiswa":
		if student.UserID != userID {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied",
			})
		}

	case "dosen_wali":
		ok, err := repository.IsStudentUnderAdvisor(
			database.PgPool,
			userID,
			student.StudentPK,
		)
		if err != nil || !ok {
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
		"data":   student,
	})
}

func GetStudentAchievements(c *fiber.Ctx) error {
	studentPK := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	student, err := repository.GetStudentByID(database.PgPool, studentPK)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	switch role {
	case "admin":
		// allowed

	case "mahasiswa":
		if student.UserID != userID {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied",
			})
		}

	case "dosen_wali":
		ok, err := repository.IsStudentUnderAdvisor(
			database.PgPool,
			userID,
			studentPK,
		)
		if err != nil || !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "not your advisee",
			})
		}

	default:
		return c.Status(403).JSON(fiber.Map{
			"error": "role not allowed",
		})
	}

	data, err := repository.GetStudentAchievements(database.PgPool, studentPK)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load achievements",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func UpdateStudentAdvisor(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    if role != "admin" {
        return c.Status(403).JSON(fiber.Map{
            "error": "admin only",
        })
    }

    studentPK := c.Params("id")

    var body model.UpdateStudentAdvisorRequest
    if err := c.BodyParser(&body); err != nil || body.AdvisorID == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "advisor_id is required",
        })
    }

    err := repository.UpdateStudentAdvisor(
        database.PgPool,
        studentPK,
        body.AdvisorID,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "failed to update advisor",
        })
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "advisor updated",
    })
}


