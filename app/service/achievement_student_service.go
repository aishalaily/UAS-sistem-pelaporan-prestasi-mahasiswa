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
			return c.Status(500).JSON(fiber.Map{"error": "failed to load advisee", "detail": err.Error()})
		}
		data, err := repository.GetAchievementsForStudents(database.PgPool, students)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load achievements", "detail": err.Error()})
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

	var req model.AchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.AchievementType == "" || req.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "achievementType & title are required"})
	}

	mongoPayload := model.AchievementMongo{
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Attachments:     []model.AchievementAttachment{}, // kosong dulu
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
			"id":       refID,
			"mongo_id": mongoID,
			"status":   "draft",
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


func VerifyAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	userID := c.Locals("user_id").(string)

	ref, err := repository.GetReferenceByID(database.PgPool, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{"error": "achievement not submitted"})
	}

	allowed, err := repository.IsStudentUnderAdvisor(
		database.PgPool,
		userID,
		ref.StudentID,
	)
	if err != nil || !allowed {
		return c.Status(403).JSON(fiber.Map{"error": "not your advisee"})
	}

	if err := repository.VerifyAchievement(database.PgPool, refID, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to verify"})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement verified",
	})
}

func RejectAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	userID := c.Locals("user_id").(string)

	var body model.RejectAchievementRequest
	if err := c.BodyParser(&body); err != nil || body.Note == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "rejection note is required",
		})
	}

	ref, err := repository.GetReferenceByID(database.PgPool, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "achievement not found",
		})
	}
	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"error": "achievement not submitted",
		})
	}

	allowed, err := repository.IsStudentUnderAdvisor(
		database.PgPool,
		userID,
		ref.StudentID,
	)
	if err != nil || !allowed {
		return c.Status(403).JSON(fiber.Map{
			"error": "not your advisee",
		})
	}

	if err := repository.RejectAchievement(
		database.PgPool,
		refID,
		body.Note,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to reject achievement",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement rejected",
	})
}

func GetAchievementHistory(c *fiber.Ctx) error {
	refID := c.Params("id")

	ref, err := repository.GetReferenceByID(database.PgPool, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	history := []model.AchievementHistory{}

	if ref.SubmittedAt != nil {
		history = append(history, model.AchievementHistory{
			Status:    "submitted",
			ChangedAt: *ref.SubmittedAt,
		})
	}

	if ref.VerifiedAt != nil {
		history = append(history, model.AchievementHistory{
			Status:    ref.Status,
			ChangedAt: *ref.VerifiedAt,
			ActorID:   ref.VerifiedBy,
			Note:      ref.RejectionNote,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   history,
	})
}

func UploadAchievementAttachment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	refID := c.Params("id")

	studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "student not found"})
	}

	ref, err := repository.GetReferenceByIDAndStudent(database.PgPool, refID, studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	form, err := c.MultipartForm()
	if err != nil || len(form.File["file"]) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "file is required"})
	}

	file := form.File["file"][0]
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	path := "./uploads/" + filename

	if err := c.SaveFile(file, path); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save file"})
	}

	attachment := model.AchievementAttachment{
		FileName:   filename,
		FileURL:    "/uploads/" + filename,
		FileType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	if err := repository.AddAchievementAttachment(ref.MongoAchievementID, attachment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to attach file"})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "file uploaded",
	})
}
