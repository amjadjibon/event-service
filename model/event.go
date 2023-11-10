package model

import (
	"time"
)

type Event struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

type EventDetail struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
	TotalWorkshop int       `json:"total_workshop"`
}
