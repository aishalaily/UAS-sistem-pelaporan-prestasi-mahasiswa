package repository

import (
	"uas-go/app/model"
	"uas-go/database"
)	

type AchievementRepository interface {
	GetStudentIDByUserID(userID string) (string, error)
	InsertAchievementMongo(data model.AchievementMongo) (string, error)
	InsertReference(studentID, mongoID string) (string, error)
	GetReferenceByID(refID string) (*model.AchievementReference, error)
	IsStudentUnderAdvisor(userID, studentID string) (bool, error)
	VerifyAchievement(refID string, points int, verifierID string) error
}

type AchievementRepositoryImpl struct{}

func (r *AchievementRepositoryImpl) GetStudentIDByUserID(userID string) (string, error) {
	return GetStudentIDByUserID(database.PgPool, userID)
}

func (r *AchievementRepositoryImpl) InsertAchievementMongo(data model.AchievementMongo) (string, error) {
	return InsertAchievementMongo(data)
}

func (r *AchievementRepositoryImpl) InsertReference(studentID, mongoID string) (string, error) {
	return InsertReference(database.PgPool, studentID, mongoID)
}

func (r *AchievementRepositoryImpl) GetReferenceByID(refID string) (*model.AchievementReference, error) {
	return GetReferenceByID(database.PgPool, refID)
}

func (r *AchievementRepositoryImpl) IsStudentUnderAdvisor(userID, studentID string) (bool, error) {
	return IsStudentUnderAdvisor(database.PgPool, userID, studentID)
}

func (r *AchievementRepositoryImpl) VerifyAchievement(refID string, points int, verifierID string) error {
	return VerifyAchievement(database.PgPool, refID, points, verifierID)
}
