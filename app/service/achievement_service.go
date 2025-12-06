package service

import (
	"uas-go/app/model"
	"uas-go/app/repository"

	"github.com/gofiber/fiber/v2"
)

func SubmitAchievement(c *fiber.Ctx) error {
	var req model.AchievementRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request",
		})
	}

	studentID := c.Locals("user_id").(string)

	mongoID, err := repository.InsertAchievementMongo(studentID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to save data to mongodb",
			"error":   err.Error(),
		})
	}

	ref, err := repository.InsertAchievementReference(studentID, mongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to save reference",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":        ref.ID,
			"mongoId":   ref.MongoAchievementID,
			"status":    ref.Status,
			"submitted": ref.SubmittedAt,
		},
	})
}

func SubmitForVerification(c *fiber.Ctx) error {

	refID := c.Params("id")
	studentID := c.Locals("user_id").(string)

	ref, err := repository.SubmitAchievementForVerification(refID, studentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "cannot submit achievement",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":        ref.ID,
			"status":    ref.Status,
			"submitted": ref.SubmittedAt,
		},
	})
}

func DeleteAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	studentID := c.Locals("user_id").(string)

	ref, err := repository.SoftDeleteAchievementReference(refID, studentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to delete achievement (must be draft)",
			"error":   err.Error(),
		})
	}

	err = repository.SoftDeleteAchievementMongo(ref.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to delete mongo data",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "achievement deleted successfully",
	})
}
