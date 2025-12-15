package model

type AchievementStatistics struct {
	ByType   map[string]int `json:"by_type"`
	ByPeriod map[string]int `json:"by_period"`
}

type TopStudent struct {
	StudentID string `json:"student_id"`
	Total     int    `json:"total"`
}

