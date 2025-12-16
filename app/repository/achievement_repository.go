package repository

import (
	"context"
	"time"

	"uas-go/app/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertReference(db *pgxpool.Pool, studentID, MongoAchievementID string) (string, error) {
	refID := uuid.NewString()
	now := time.Now()

	_, err := db.Exec(context.Background(), `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, 'draft', $4, $5, false)
	`, refID, studentID, MongoAchievementID, now, now)

	if err != nil {
		return "", err
	}
	return refID, nil
}

func GetReferenceByID(db *pgxpool.Pool, id string) (*model.AchievementReference, error) {
	row := db.QueryRow(context.Background(), `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at, is_deleted
		FROM achievement_references
		WHERE id = $1 
		AND is_deleted = false
		LIMIT 1
	`, id)

	var ref model.AchievementReference
	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt, &ref.IsDeleted,
	)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func GetReferenceByIDAndStudent(db *pgxpool.Pool, id, studentID string) (*model.AchievementReference, error) {
	row := db.QueryRow(context.Background(), `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at, is_deleted
		FROM achievement_references
		WHERE id = $1 
		AND student_id = $2 
		AND is_deleted = false
		LIMIT 1
	`, id, studentID)

	var ref model.AchievementReference
	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt, &ref.IsDeleted,
	)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func SubmitForVerification(db *pgxpool.Pool, id string) error {
	_, err := db.Exec(context.Background(), `
		UPDATE achievement_references
		SET status = 'submitted', 
			submitted_at = NOW(), 
			updated_at = NOW()
		WHERE id = $1 
		AND status = 'draft'
	`, id)
	return err
}

func SoftDeleteReference(db *pgxpool.Pool, id, studentID string) error {
	_, err := db.Exec(context.Background(), `
		UPDATE achievement_references
		SET is_deleted = true, 
			updated_at = NOW()
		WHERE id = $1 
		AND student_id = $2
	`, id, studentID)
	return err
}

func GetAchievementsByStudent(db *pgxpool.Pool, studentID string) ([]model.AchievementReference, error) {
	rows, err := db.Query(context.Background(), `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.created_at
		FROM achievement_references ar
		WHERE ar.student_id = $1 
		AND ar.is_deleted = false
		ORDER BY ar.created_at DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.AchievementReference
	for rows.Next() {
		var a model.AchievementReference
		if err := rows.Scan(&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, nil
}

func GetAchievementsForStudents(db *pgxpool.Pool, student []string) ([]model.AchievementReference, error) {
	if len(student) == 0 {
		return []model.AchievementReference{}, nil
	}

	rows, err := db.Query(context.Background(), `
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE student_id = ANY($1)
		AND is_deleted = false
		AND status != 'draft'
		ORDER BY created_at DESC
	`, student)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.AchievementReference
	for rows.Next() {
		var a model.AchievementReference
		if err := rows.Scan(&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, nil
}

func GetAllAchievements(db *pgxpool.Pool) ([]model.AchievementReference, error) {
	rows, err := db.Query(context.Background(), `
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE is_deleted = false
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.AchievementReference
	for rows.Next() {
		var a model.AchievementReference
		if err := rows.Scan(&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, nil
}

func VerifyAchievement(db *pgxpool.Pool, refID string, points int, verifierID string) error {
	_, err := db.Exec(context.Background(), `
		UPDATE achievement_references
		SET status = 'verified',
		    points = $2,
		    verified_at = NOW(),
		    verified_by = $3,
		    updated_at = NOW()
		WHERE id = $1
		  AND status = 'submitted'
	`, refID, points, verifierID)

	return err
}

func RejectAchievement(db *pgxpool.Pool, refID, note string) error {
	_, err := db.Exec(context.Background(), `
		UPDATE achievement_references
		SET status = 'rejected',
		    rejection_note = $2,
		    updated_at = NOW()
		WHERE id = $1 AND status = 'submitted'
	`, refID, note)

	return err
}
