package test

import (
	"bytes"
	"net/http/httptest"
	"testing"
	
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"uas-go/app/model"
	"uas-go/app/service"
	"uas-go/app/repository/mocks"
)

func TestLoginSuccess(t *testing.T) {
	app := fiber.New()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mockRepo := &mocks.UserRepositoryMock{
		GetUserByUsernameFn: func(username string) (*model.User, string, error) {
			return &model.User{
				ID:       "user-1",
				Username: "admin",
				RoleID:   "role-1",
				FullName: "Admin Test",
			}, string(hashed), nil
		},
		GetRoleNameFn: func(roleID string) (string, error) {
			return "admin", nil
		},
		GetPermissionsFn: func(roleID string) ([]string, error) {
			return []string{"achievement.read"}, nil
		},
	}

	authService := service.NewAuthService(mockRepo)
	app.Post("/login", authService.Login)

	body := []byte(`{"username":"admin","password":"password"}`)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLoginWrongPassword(t *testing.T) {
	app := fiber.New()

	mockRepo := &mocks.UserRepositoryMock{
		GetUserByUsernameFn: func(username string) (*model.User, string, error) {
			return &model.User{
				ID:       "user-1",
				Username: "admin",
			}, "wrong-hash", nil
		},
	}

	authService := service.NewAuthService(mockRepo)

	app.Post("/login", authService.Login)

	body := []byte(`{"username":"admin","password":"wrong"}`)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}