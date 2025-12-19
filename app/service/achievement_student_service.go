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

var AchievementRepo repository.AchievementRepository = &repository.AchievementRepositoryImpl{}

// GetAchievements godoc
// @Summary Get achievements
// @Description Get achievements list based on role (Mahasiswa, Dosen Wali, Admin)
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements [get]
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

// GetAchievementDetail godoc
// @Summary Get achievement detail
// @Description Get achievement detail by reference ID
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /achievements/{id} [get]
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

// SubmitAchievement godoc
// @Summary Create achievement
// @Description Create new achievement (Mahasiswa only)
// @Tags Achievements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.AchievementRequest true "Achievement payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /achievements [post]
func SubmitAchievement(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	studentID, err := AchievementRepo.GetStudentIDByUserID(userID)
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

	mongoID, err := AchievementRepo.InsertAchievementMongo(mongoPayload)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save mongo"})
	}

	refID, err := AchievementRepo.InsertReference(studentID, mongoID)
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

// SubmitForVerification godoc
// @Summary Submit achievement
// @Description Submit achievement for verification (Mahasiswa only)
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]string
// @Router /achievements/{id}/submit [post]
func SubmitForVerification(c *fiber.Ctx) error {
	refID := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "dosen_wali" {
		return c.Status(403).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

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

// UpdateAchievement godoc
// @Summary Update achievement
// @Description Update draft achievement (Mahasiswa only)
// @Tags Achievements
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /achievements/{id} [put]
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

// DeleteAchievement godoc
// @Summary Delete achievement
// @Description Delete draft achievement (Mahasiswa only)
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /achievements/{id} [delete]
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

// VerifyAchievement godoc
// @Summary Verify achievement
// @Description Verify achievement and assign points (Dosen Wali only)
// @Tags Achievements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Param body body model.VerifyAchievementRequest true "Verification payload"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/verify [post]
func VerifyAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	userID := c.Locals("user_id").(string)

	var body model.VerifyAchievementRequest
	if err := c.BodyParser(&body); err != nil || body.Points <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "points is required and must be > 0",
		})
	}

	ref, err := AchievementRepo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{"error": "achievement not submitted"})
	}

	allowed, err := AchievementRepo.IsStudentUnderAdvisor(
		userID,
		ref.StudentID,
	)
	if err != nil || !allowed {
		return c.Status(403).JSON(fiber.Map{"error": "not your advisee"})
	}

	if err := AchievementRepo.VerifyAchievement(
		refID,
		body.Points,
		userID,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to verify"})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement verified",
		"points":  body.Points,
	})
}

// RejectedAchievement godoc
// @Summary Reject achievement
// @Description Reject achievement (Dosen Wali only)
// @Tags Achievements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Param body body model.RejectAchievementRequest true "Rejection payload"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/reject [post]
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

// GetAchievementHistory godoc
// @Summary Get achievement history
// @Description Get achievement status history
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {array} model.AchievementHistory
// @Router /achievements/{id}/history [get]
func GetAchievementHistory(c *fiber.Ctx) error {
	refID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	ref, err := repository.GetReferenceByID(database.PgPool, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "achievement not found",
		})
	}

	switch role {
	case "mahasiswa":
		studentID, err := repository.GetStudentIDByUserID(database.PgPool, userID)
		if err != nil || studentID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{"error": "access denied"})
		}

	case "dosen_wali":
		allowed, err := repository.IsStudentUnderAdvisor(
			database.PgPool,
			userID,
			ref.StudentID,
		)
		if err != nil || !allowed {
			return c.Status(403).JSON(fiber.Map{"error": "not your advisee"})
		}

	case "admin":
		// allowed

	default:
		return c.Status(403).JSON(fiber.Map{"error": "role not allowed"})
	}

	var history []model.AchievementHistory

	history = append(history, model.AchievementHistory{
		Status:    "draft",
		ChangedAt: ref.CreatedAt,
	})

	if ref.SubmittedAt != nil {
		history = append(history, model.AchievementHistory{
			Status:    "submitted",
			ChangedAt: *ref.SubmittedAt,
		})
	}

	getActorName := func(actorID *string) *string {
		if actorID == nil {
			return nil
		}
		user, err := repository.GetUserByID(*actorID)
		if err != nil {
			return nil
		}
		return &user.FullName
	}

	if ref.VerifiedAt != nil && ref.Status == "verified" {
		history = append(history, model.AchievementHistory{
			Status:    "verified",
			ChangedAt: *ref.VerifiedAt,
			ActorID:   ref.VerifiedBy,
			ActorName: getActorName(ref.VerifiedBy),
		})
	}

	if ref.Status == "rejected" && ref.RejectionNote != nil {
		history = append(history, model.AchievementHistory{
			Status:    "rejected",
			ChangedAt: ref.UpdatedAt,
			ActorID:   ref.VerifiedBy,
			ActorName: getActorName(ref.VerifiedBy),
			Note:      ref.RejectionNote,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   history,
	})
}

// UploadAchievementAttachment godoc
// @Summary Upload achievement attachment
// @Description Upload attachment for achievement (Mahasiswa only)
// @Tags Achievements
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Param file formData file true "Attachment file"
// @Success 200 {object} map[string]string
// @Router /achievements/{id}/attachments [post]
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
