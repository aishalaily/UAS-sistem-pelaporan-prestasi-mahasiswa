package model

import "time"

type AchievementReference struct {
	ID                 string    `json:"id" bson:"id"`
	StudentID          string    `json:"student_id" bson:"student_id"`
	MongoAchievementID string    `json:"mongo_achievement_id" bson:"mongo_achievement_id"`
	Status             string    `json:"status"` // draft, submitted, verified, rejected

	SubmittedAt   *time.Time `json:"submitted_at"`
	VerifiedAt    *time.Time `json:"verified_at"`
	VerifiedBy    *string    `json:"verified_by"`
	RejectionNote *string    `json:"rejection_note"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

type AchievementAttachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

type AchievementMongo struct {
	ID              interface{}            `bson:"_id" json:"id"`
	StudentID       string                 `bson:"studentId" json:"studentId"`
	AchievementType string                 `bson:"achievementType" json:"achievementType"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`

	Details     map[string]interface{} `bson:"details" json:"details"`
	Tags        []string               `bson:"tags" json:"tags"`
	Attachments []AchievementAttachment `bson:"attachments" json:"attachments"`

	Points    int       `bson:"points" json:"points"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type AchievementRequest struct {
	AchievementType string                 `json:"achievementType"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
}

type AchievementHistory struct {
	Status    string     `json:"status"`
	ChangedAt time.Time  `json:"changed_at"`
	ActorID   *string    `json:"actor_id"`
	Note      *string    `json:"note"`
}
