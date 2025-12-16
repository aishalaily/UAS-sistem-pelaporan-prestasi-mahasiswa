package model

type AchievementStatistics struct {
	ByType      map[string]int `json:"by_type"`
	ByPeriod    map[string]int `json:"by_period"`
	TopStudents []TopStudent   `json:"top_students,omitempty"`
	Competition map[string]int `json:"competition_distribution,omitempty"`
}

type TopStudent struct {
	StudentID string `json:"student_id"`
	Total     int    `json:"total"`
}