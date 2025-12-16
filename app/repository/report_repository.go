package repository

import (
	"context"

	"uas-go/app/model"
	"uas-go/database"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAchievementStatsAdmin() (map[string]int, error) {
	rows, err := database.PgPool.Query(context.Background(), `
		SELECT EXTRACT(YEAR FROM created_at)::TEXT, COUNT(*)
		FROM achievement_references
		WHERE status = 'verified'
		GROUP BY 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var year string
		var count int

		if err := rows.Scan(&year, &count); err != nil {
			return nil, err
		}
		result[year] = count
	}

	return result, nil
}

func GetAchievementStatsStudent(studentID string) (map[string]int, error) {
	rows, err := database.PgPool.Query(context.Background(), `
		SELECT EXTRACT(YEAR FROM created_at)::TEXT, COUNT(*)
		FROM achievement_references
		WHERE status = 'verified'
		AND student_id = $1
		GROUP BY 1
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var year string
		var count int

		if err := rows.Scan(&year, &count); err != nil {
			return nil, err
		}
		result[year] = count
	}


	return result, nil
}

func GetAchievementStatsForStudents(studentIDs []string) (map[string]int, error) {
	rows, err := database.PgPool.Query(context.Background(), `
		SELECT EXTRACT(YEAR FROM created_at)::TEXT, COUNT(*)
		FROM achievement_references
		WHERE status = 'verified'
		AND student_id = ANY($1)
		GROUP BY 1
	`, studentIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var year string
		var count int

		if err := rows.Scan(&year, &count); err != nil {
			return nil, err
		}
		result[year] = count
	}

	return result, nil
}

func GetAchievementTypeStatsMongo(studentIDs []string) (map[string]int, error) {
	collection := database.MongoDB.Collection("achievements")

	filter := bson.M{}
	if len(studentIDs) > 0 {
		filter["studentId"] = bson.M{"$in": studentIDs}
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	result := make(map[string]int)

	for cursor.Next(context.Background()) {
		var doc struct {
			AchievementType string `bson:"achievementType"`
		}

		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		if doc.AchievementType != "" {
			result[doc.AchievementType]++
		}
	}

	return result, nil
}

func GetTopStudents(limit int) ([]model.TopStudent, error) {
	rows, err := database.PgPool.Query(context.Background(), `
		SELECT student_id, COUNT(*) AS total
		FROM achievement_references
		WHERE status = 'verified'
		GROUP BY student_id
		ORDER BY total DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.TopStudent
	for rows.Next() {
		var s model.TopStudent
		if err := rows.Scan(&s.StudentID, &s.Total); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func GetCompetitionLevelDistribution(studentIDs []string) (map[string]int, error) {
	collection := database.MongoDB.Collection("achievements")

	filter := bson.M{}
	if len(studentIDs) > 0 {
		filter["studentId"] = bson.M{"$in": studentIDs}
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	result := make(map[string]int)

	for cursor.Next(context.Background()) {
		var doc struct {
			Details map[string]interface{} `bson:"details"`
		}

		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		// ambil level kompetisi
		if level, ok := doc.Details["level"].(string); ok {
			result[level]++
		}
	}

	return result, nil
}
