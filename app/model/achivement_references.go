package model

import "time"

type AchievementReference struct {
	ID                string    `json:"id"`
	StudentID         string    `json:"student_id"`
	MongoAchievementID string   `json:"mongo_achievement_id"`
	Status            string    `json:"status"` // draft, submitted, verified, rejected
	SubmittedAt       time.Time `json:"submitted_at"`
	VerifiedAt        time.Time `json:"verified_at"`
	VerifiedBy        string    `json:"verified_by"`
	RejectionNote     string    `json:"rejection_note"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	IsDeleted         bool      `json:"is_deleted"`
}

type AchievementRequest struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Category    string   `json:"category"`
    EventDate   string   `json:"event_date"`
    Documents   []string `json:"documents"`
}

type AchievementRef struct {
    ID        string    `json:"id"`
    StudentID string    `json:"student_id"`
    Title     string    `json:"title"`
    Category  string    `json:"category"`
    Status    string    `json:"status"`
    MongoID   string    `json:"mongo_id"`
    CreatedAt time.Time `json:"created_at"`
}
