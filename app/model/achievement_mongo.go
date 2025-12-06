package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	StudentID      string             `bson:"studentId"`
	StudentSnapshot struct {
		StudentID    string `bson:"student_id"`
		FullName     string `bson:"fullName"`
		ProgramStudy string `bson:"programStudy"`
	} `bson:"studentSnapshot"`
	Title       string                 `bson:"title"`
	Description string                 `bson:"description"`
	Type        string                 `bson:"achievementType"`
	Details     map[string]interface{} `bson:"details"`
	Attachments []Attachment           `bson:"attachments"`
	Points      int                    `bson:"points"`
	CreatedAt   time.Time              `bson:"createdAt"`
	UpdatedAt   time.Time              `bson:"updatedAt"`
}

type Attachment struct {
	FileName  string    `bson:"fileName"`
	FileURL   string    `bson:"fileUrl"`
	FileType  string    `bson:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt"`
}

