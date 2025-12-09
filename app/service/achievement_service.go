package service

import (
	"encoding/json"
	"fmt"
	"time"

	"uas-go/app/model"
	"uas-go/app/repository"
	"uas-go/database"

	"github.com/gofiber/fiber/v2"
)

func GetAchievements(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	switch role {
	case "mahasiswa":
		studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "student profile not found"})
		}
		data, err := repository.GetAchievementsByStudent(database.PgPool, studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load achievements"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": data})

	case "dosen_wali":
		students, err := repository.GetStudentsUnderAdvisor(database.PgPool, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load advisee"})
		}
		data, err := repository.GetAchievementsForStudents(database.PgPool, students)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load achievements"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": data})

	case "admin":
		data, err := repository.GetAllAchievements(database.PgPool)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load achievements"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": data})

	default:
		return c.Status(403).JSON(fiber.Map{"error": "role not allowed"})
	}
}

func GetAchievementDetail(c *fiber.Ctx) error {
	refID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := repository.GetReferenceByID(database.PgPool, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	switch role {
	case "mahasiswa":
		sid, err := repository.GetStudentIDByUserID(database.PgPool, userID)
		if err != nil || sid != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "access denied"})
		}
	case "dosen_wali":
		isChild, err := repository.IsStudentUnderAdvisor(database.PgPool, userID, ref.StudentID)
		if err != nil || !isChild {
			return c.Status(403).JSON(fiber.Map{"error": "not your advisee"})
		}
	case "admin":
	default:
		return c.Status(403).JSON(fiber.Map{"error": "role not allowed"})
	}

	mongoData, err := repository.GetAchievementMongo(ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "mongo data not found"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"reference": ref,
			"details":   mongoData,
		},
	})
}

func SubmitAchievement(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "student profile not found"})
	}

	achievementType := c.FormValue("achievementType")
	title := c.FormValue("title")
	description := c.FormValue("description")
	if achievementType == "" || title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "achievementType & title are required"})
	}

	var details map[string]interface{}
	if err := json.Unmarshal([]byte(c.FormValue("details")), &details); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid details JSON"})
	}

	var tags []string
	_ = json.Unmarshal([]byte(c.FormValue("tags")), &tags)

	form, _ := c.MultipartForm()
	var attachments []model.AchievementAttachment
	if form != nil {
		for _, f := range form.File["documents"] {
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), f.Filename)
			savePath := "./uploads/" + filename
			if err := c.SaveFile(f, savePath); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "failed to save file"})
			}
			attachments = append(attachments, model.AchievementAttachment{
				FileName:   filename,
				FileURL:    "/uploads/" + filename,
				FileType:   f.Header.Get("Content-Type"),
				UploadedAt: time.Now(),
			})
		}
	}

	mongoPayload := model.AchievementMongo{
		StudentID:       studentID,
		AchievementType: achievementType,
		Title:           title,
		Description:     description,
		Details:         details,
		Tags:            tags,
		Attachments:     attachments,
		Points:          0,
	}

	mongoID, err := repository.InsertAchievementMongo(mongoPayload)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save mongo"})
	}

	refID, err := repository.InsertReference(database.PgPool, studentID, mongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save reference"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":         refID,
			"mongo_id":   mongoID,
			"student_id": studentID,
			"status":     "draft",
		},
	})
}

func SubmitForVerification(c *fiber.Ctx) error {
	refID := c.Params("id")
	userID := c.Locals("user_id").(string)

	studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "student profile not found"})
	}

	ref, err := repository.GetReferenceByIDAndStudent(database.PgPool, refID, studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft achievement can be submitted"})
	}

	if err := repository.SubmitForVerification(database.PgPool, refID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update status"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "submitted"})
}

func UpdateAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	userID := c.Locals("user_id").(string)

	studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "student profile not found"})
	}

	ref, err := repository.GetReferenceByIDAndStudent(database.PgPool, refID, studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft achievement can be updated"})
	}

	oldData, err := repository.GetAchievementMongo(ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "mongo data not found"})
	}

	if c.FormValue("achievementType") != "" {
		oldData.AchievementType = c.FormValue("achievementType")
	}
	if c.FormValue("title") != "" {
		oldData.Title = c.FormValue("title")
	}
	if c.FormValue("description") != "" {
		oldData.Description = c.FormValue("description")
	}
	if c.FormValue("details") != "" {
		var details map[string]interface{}
		if err := json.Unmarshal([]byte(c.FormValue("details")), &details); err == nil {
			oldData.Details = details
		}
	}
	if c.FormValue("tags") != "" {
		var tags []string
		if err := json.Unmarshal([]byte(c.FormValue("tags")), &tags); err == nil {
			oldData.Tags = tags
		}
	}

	form, _ := c.MultipartForm()
	if form != nil {
		for _, f := range form.File["documents"] {
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), f.Filename)
			path := "./uploads/" + filename
			c.SaveFile(f, path)
			oldData.Attachments = append(oldData.Attachments, model.AchievementAttachment{
				FileName:   filename,
				FileURL:    "/uploads/" + filename,
				FileType:    f.Header.Get("Content-Type"),
				UploadedAt: time.Now(),
			})
		}
	}

	if err := repository.UpdateAchievementMongo(ref.MongoAchievementID, *oldData); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "updated"})
}

func DeleteAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	userID := c.Locals("user_id").(string)

	studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "student profile not found"})
	}

	ref, err := repository.GetReferenceByIDAndStudent(database.PgPool, refID, studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft achievement can be deleted"})
	}

	if err := repository.SoftDeleteReference(database.PgPool, refID, studentID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "deleted"})
}
