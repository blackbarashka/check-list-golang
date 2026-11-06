package tasks

import "time"

type Task struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Status        bool       `json:"status"`
	AddTime       time.Time  `json:"add_time"`
	CompletedTime *time.Time `json:"completed_by"`
}

func NewTask(title, description string) *Task {
	return &Task{
		Title:       title,
		Description: description,
		AddTime:     time.Now(),
		Status:      false,
	}
}
