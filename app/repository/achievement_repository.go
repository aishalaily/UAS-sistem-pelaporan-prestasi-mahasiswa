package repository

import (
	"context"
	"time"
	"uas-go/app/model"
	"uas-go/database"
	"database/sql"

	"github.com/google/uuid"
)

func InsertAchievementReference(studentID, mongoID string) (*model.AchievementReference, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO achievement_references (
			id, student_id, mongo_achievement_id, status,
			submitted_at, created_at, updated_at, is_deleted
		)
		VALUES ($1, $2, $3, 'draft', $4, $5, $6, false)
		RETURNING id, student_id, mongo_achievement_id, status,
		          submitted_at, verified_at, verified_by,
				  rejection_note, created_at, updated_at, is_deleted
	`

	row := database.PgPool.QueryRow(context.Background(),
		query, id, studentID, mongoID, now, now, now)

	var ref model.AchievementReference

	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
		&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt, &ref.IsDeleted,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

func GetAchievementReference(db *sql.DB, id string) (*model.AchievementReference, error) {
    var ref model.AchievementReference

    row := db.QueryRow(`
        SELECT id, student_id, mongo_achievement_id, status,
               submitted_at, verified_at, verified_by,
               rejection_note, created_at, updated_at, is_deleted
        FROM achievement_references
        WHERE id = $1 AND is_deleted = false
        LIMIT 1
    `, id)

    err := row.Scan(
        &ref.ID,
        &ref.StudentID,
        &ref.MongoAchievementID,
        &ref.Status,
        &ref.SubmittedAt,
        &ref.VerifiedAt,
        &ref.VerifiedBy,
        &ref.RejectionNote,
        &ref.CreatedAt,
        &ref.UpdatedAt,
        &ref.IsDeleted,
    )

    if err != nil {
        return nil, err
    }

    return &ref, nil
}


func SubmitAchievementForVerification(refID string, studentID string) (*model.AchievementReference, error) {
	now := time.Now()

	query := `
		UPDATE achievement_references
		SET status = 'submitted',
		    submitted_at = $1,
		    updated_at = $1
		WHERE id = $2 AND student_id = $3 AND status = 'draft'
		RETURNING id, student_id, mongo_achievement_id, status,
		          submitted_at, verified_at, verified_by,
		          rejection_note, created_at, updated_at, is_deleted
	`

	row := database.PgPool.QueryRow(context.Background(), query, now, refID, studentID)

	var ref model.AchievementReference

	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
		&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt, &ref.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &ref, nil
}

func UpdateAchievementStatus(db *sql.DB, id string, newStatus string) error {
    _, err := db.Exec(`
        UPDATE achievement_references
        SET status = $1, submitted_at = NOW(), updated_at = NOW()
        WHERE id = $2
    `, newStatus, id)

    return err
}

func SoftDeleteAchievementReference(refID string, studentID string) (*model.AchievementReference, error) {
	now := time.Now()

	query := `
		UPDATE achievement_references
		SET is_deleted = true,
		    updated_at = $1
		WHERE id = $2 AND student_id = $3 AND status = 'draft'
		RETURNING id, student_id, mongo_achievement_id, status,
		          submitted_at, verified_at, verified_by,
		          rejection_note, created_at, updated_at, is_deleted
	`

	row := database.PgPool.QueryRow(context.Background(), query, now, refID, studentID)

	var ref model.AchievementReference

	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
		&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt, &ref.IsDeleted,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}


