package models

import "time"

type Lesson struct {
	ID        string `json:"id"`
	Subject   string `json:"subject"`
	Date      string `json:"date"`       // формат YYYY-MM-DD
	StartTime string `json:"start_time"` // HH:MM
	EndTime   string `json:"end_time"`   // HH:MM
}

// Проверка пересечения времён
func (l *Lesson) OverlapsWith(other *Lesson) bool {
	if l.Date != other.Date {
		return false
	}

	start1, _ := time.Parse("15:04", l.StartTime)
	end1, _ := time.Parse("15:04", l.EndTime)
	start2, _ := time.Parse("15:04", other.StartTime)
	end2, _ := time.Parse("15:04", other.EndTime)

	return start1.Before(end2) && start2.Before(end1)
}
