package mocks

import (
	"uas-go/app/model"
)

type AchievementRepositoryMock struct {
	GetStudentErr error
	GetRefFn              func(string) (*model.AchievementReference, error)
	IsAdvisorFn           func(string, string) (bool, error)
	VerifyAchievementFn   func(string, int, string) error
}

func (m *AchievementRepositoryMock) GetStudentIDByUserID(userID string) (string, error) {
	if m.GetStudentErr != nil {
		return "", m.GetStudentErr
	}
	return "student-123", nil
}

func (m *AchievementRepositoryMock) InsertAchievementMongo(data model.AchievementMongo) (string, error) {
	return "mongo-123", nil
}

func (m *AchievementRepositoryMock) InsertReference(studentID, mongoID string) (string, error) {
	return "ref-123", nil
}

func (m *AchievementRepositoryMock) GetReferenceByID(id string) (*model.AchievementReference, error) {
	return m.GetRefFn(id)
}

func (m *AchievementRepositoryMock) IsStudentUnderAdvisor(uid, sid string) (bool, error) {
	return m.IsAdvisorFn(uid, sid)
}

func (m *AchievementRepositoryMock) VerifyAchievement(id string, p int, uid string) error {
	return m.VerifyAchievementFn(id, p, uid)
}