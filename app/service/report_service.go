package service

import (
	"uas-go/app/repository"
	"uas-go/database"
	"uas-go/app/model"

	"github.com/gofiber/fiber/v2"
)

func GetAchievementStatistics(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	var result model.AchievementStatistics

	switch role {

	case "admin":
		var err error

		result.ByPeriod, err = repository.GetAchievementStatsAdmin()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load period statistics"})
		}

		result.ByType, err = repository.GetAchievementTypeStatsMongo(nil)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load type statistics"})
		}

		result.TopStudents, err = repository.GetTopStudents(5)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load top students"})
		}

	case "mahasiswa":
		studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "student not found"})
		}

		result.ByPeriod, err = repository.GetAchievementStatsStudent(studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load period statistics"})
		}

		result.ByType, err = repository.GetAchievementTypeStatsMongo([]string{studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load type statistics"})
		}

		result.Competition, err = repository.GetCompetitionLevelDistribution([]string{studentID})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load competition distribution"})
		}

	case "dosen_wali":
		studentIDs, err := repository.GetStudentsUnderAdvisor(database.PgPool, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load advisees"})
		}

		result.ByPeriod, err = repository.GetAchievementStatsForStudents(studentIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load period statistics"})
		}

		result.ByType, err = repository.GetAchievementTypeStatsMongo(studentIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load type statistics"})
		}

		result.Competition, err = repository.GetCompetitionLevelDistribution(studentIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load competition distribution"})
		}

	default:
		return c.Status(403).JSON(fiber.Map{"error": "role not allowed"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func GetStudentReport(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	studentID := c.Params("id")

	switch role {

	case "admin":

	case "mahasiswa":
		ownStudentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
		if err != nil || ownStudentID != studentID {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied",
			})
		}

	case "dosen_wali":
		isAllowed, err := repository.IsStudentUnderAdvisor(
			database.PgPool,
			userID,
			studentID,
		)
		if err != nil || !isAllowed {
			return c.Status(403).JSON(fiber.Map{
				"error": "not your advisee",
			})
		}

	default:
		return c.Status(403).JSON(fiber.Map{
			"error": "role not allowed",
		})
	}

	byPeriod, err := repository.GetAchievementStatsStudent(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load statistics",
		})
	}

	byType, err := repository.GetAchievementTypeStatsMongo([]string{studentID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load statistics",
		})
	}

	competition, err := repository.GetCompetitionLevelDistribution([]string{studentID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load competition distribution",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"by_period": byPeriod,
			"by_type": byType,
			"top_students": []interface{}{},
			"competition_distribution": competition,
		},
	})
}
