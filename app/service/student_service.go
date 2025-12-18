package service

import (
	"uas-go/app/repository"
	"uas-go/database"
	"uas-go/app/model"

	"github.com/gofiber/fiber/v2"
)

// GetStudents godoc
// @Summary Get all students
// @Description Admin get all students
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /students [get]
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

// GetStudentByID godoc
// @Summary Get student detail
// @Description Get student detail by ID (Admin, Mahasiswa own, Dosen Wali advisee)
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID (PK)"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /students/{id} [get]
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

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Get achievements of a student
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID (PK)"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /students/{id}/achievements [get]
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

// UpdateStudentAdvisor godoc
// @Summary Update student advisor
// @Description Assign or update advisor for student (Admin only)
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID (PK)"
// @Param body body model.UpdateStudentAdvisorRequest true "Advisor payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /students/{id}/advisor [put]
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