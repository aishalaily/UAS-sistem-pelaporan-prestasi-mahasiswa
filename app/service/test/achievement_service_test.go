package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"uas-go/app/model"
	mockRepo "uas-go/app/repository/mocks"
	"uas-go/app/service"
)

func TestSubmitAchievementSuccess(t *testing.T) {
	app := fiber.New()

	mock := &mockRepo.AchievementRepositoryMock{}
	service.AchievementRepo = mock

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return service.SubmitAchievement(c)
	})

	payload := model.AchievementRequest{
		AchievementType: "competition",
		Title:           "Juara 1",
		Description:     "Lomba nasional",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/achievements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestVerifyAchievementSuccess(t *testing.T) {
	// === MOCK ===
	mock := &mockRepo.AchievementRepositoryMock{
		GetRefFn: func(id string) (*model.AchievementReference, error) {
			return &model.AchievementReference{
				ID:        id,
				StudentID: "student-123",
				Status:    "submitted",
			}, nil
		},
		IsAdvisorFn: func(uid, sid string) (bool, error) {
			return true, nil
		},
		VerifyAchievementFn: func(id string, p int, uid string) error {
			return nil
		},
	}

	service.AchievementRepo = mock

	app := fiber.New()
	app.Post("/verify/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "lecturer-1")
		return service.VerifyAchievement(c)
	})

	body := map[string]int{"points": 10}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/verify/ach-1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
