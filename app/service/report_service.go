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

	var (
		byPeriod 	map[string]int
		byType  	map[string]int
		topStudents []model.TopStudent
		competition map[string]int
		err     	error
	)

	switch role {

	case "admin":
		byPeriod, err = repository.GetAchievementStatsAdmin()
		if err != nil {
			break
		}

		byType, err = repository.GetAchievementTypeStatsAdmin()
		if err != nil {
			break
		}

		topStudents, err = repository.GetTopStudents(5)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load top students"})
		}

		competition, err = repository.GetCompetitionLevelDistribution(nil)
		if err != nil {
			break
		}

	case "mahasiswa":
		studentID, err2 := repository.GetStudentIDByUserID(database.PgPool, userID)
		if err2 != nil {
			return c.Status(403).JSON(fiber.Map{"error": "student not found"})
		}

		byPeriod, err = repository.GetAchievementStatsStudent(studentID)
		if err != nil {
			break
		}
		byType, err = repository.GetAchievementTypeStatsStudent(studentID)
		if err != nil {
			break
		}
		competition, err = repository.GetCompetitionLevelDistribution([]string{studentID})
		if err != nil {
			break
		}

	case "dosen_wali":
		studentIDs, err2 := repository.GetStudentsUnderAdvisor(database.PgPool, userID)
		if err2 != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load advisees"})
		}

		byPeriod, err = repository.GetAchievementStatsForStudents(studentIDs)
		if err != nil {
			break
		}
		byType, err = repository.GetAchievementTypeStatsForStudents(studentIDs)
		if err != nil {
			break
		}
		competition, err = repository.GetCompetitionLevelDistribution(studentIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load competition distribution"})
		}

	default:
		return c.Status(403).JSON(fiber.Map{"error": "role not allowed"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load statistics"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"by_period": byPeriod,
			"by_type": byType,
			"top_students": topStudents,
			"competition_distribution": competition,
		},
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

	byType, err := repository.GetAchievementTypeStatsStudent(studentID)
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
