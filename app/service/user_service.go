package service

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"uas-go/app/model"
	"uas-go/app/repository"
	"uas-go/utils"
)

func CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid request",
		})
	}

	roleName := strings.ToLower(req.Role)

	roleID, err := repository.GetRoleIDByName(roleName)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid role",
		})
	}

	passHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to hash password",
		})
	}

	userID := uuid.New().String()

	user := model.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   roleID,
	}

	if repository.IsUsernameExists(req.Username) {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Username already exists",
		})
	}

	_, err = repository.CreateUser(user, passHash)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to create user",
		})
	}

	if roleName == "mahasiswa" {
		student := model.Student{
			ID:           uuid.New().String(),
			UserID:       userID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
			AdvisorID:    req.AdvisorID,
		}
		repository.CreateStudent(student)
	}

	if roleName == "dosen_wali" {
		lect := model.Lecturer{
			ID:         uuid.New().String(),
			UserID:     userID,
			NIDN:       req.NIDN,
			Department: req.Department,
		}
		repository.CreateLecturer(lect)
	}
	

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "User created successfully",
		"data": fiber.Map{
			"user_id": userID,
		},
	})
}
