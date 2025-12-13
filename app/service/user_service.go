package service

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"uas-go/app/model"
	"uas-go/app/repository"
	"uas-go/utils"
)

func GetUsers(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load users",
		})
	}

	var result []fiber.Map
	for _, u := range users {
		roleName, _ := repository.GetRoleName(u.RoleID)

		result = append(result, fiber.Map{
			"id":         u.ID,
			"username":   u.Username,
			"email":      u.Email,
			"full_name":  u.FullName,
			"role":       roleName,
			"is_active":  u.IsActive,
			"created_at": u.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

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

func GetUserDetail(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	userID := c.Params("id")

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	roleName, _ := repository.GetRoleName(user.RoleID)

	resp := fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"full_name":  user.FullName,
		"role":       roleName,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
	}

	if roleName == "mahasiswa" {
		student, err := repository.GetStudentByUserID(user.ID)
		if err == nil {
			resp["student"] = student
		}
	}

	if roleName == "dosen_wali" {
		lecturer, err := repository.GetLecturerByUserID(user.ID)
		if err == nil {
			resp["lecturer"] = lecturer
		}
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	userID := c.Params("id")

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	var body struct {
		Username  *string `json:"username"`
		Email     *string `json:"email"`
		FullName  *string `json:"full_name"`
		IsActive  *bool   `json:"is_active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if body.Username != nil {
		user.Username = *body.Username
	}
	if body.Email != nil {
		user.Email = *body.Email
	}
	if body.FullName != nil {
		user.FullName = *body.FullName
	}
	if body.IsActive != nil {
		user.IsActive = *body.IsActive
	}

	if err := repository.UpdateUser(*user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "user updated",
	})
}

func UpdateUserRole(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	userID := c.Params("id")

	var body struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil || body.Role == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "role is required",
		})
	}

	newRole := strings.ToLower(body.Role)

	roleID, err := repository.GetRoleIDByName(newRole)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid role",
		})
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if err := repository.UpdateUserRole(user.ID, roleID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update role",
		})
	}

	if newRole == "mahasiswa" {
		_, err := repository.GetStudentByUserID(user.ID)
		if err != nil {
			student := model.Student{
				ID:           uuid.New().String(),
				UserID:       user.ID,
				StudentID:    "",
				ProgramStudy: "",
				AcademicYear: "",
				AdvisorID:    "",
			}
			repository.CreateStudent(student)
		}
	}

	if newRole == "dosen_wali" {
		_, err := repository.GetLecturerByUserID(user.ID)
		if err != nil {
			lect := model.Lecturer{
				ID:         uuid.New().String(),
				UserID:     user.ID,
				NIDN:       "",
				Department: "",
			}
			repository.CreateLecturer(lect)
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "role updated",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin only",
		})
	}

	userID := c.Params("id")

	_, err := repository.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if err := repository.DeactivateUser(userID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "user deactivated",
	})
}
