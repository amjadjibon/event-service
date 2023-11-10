package model

import (
	"time"
)

type Workshop struct {
	ID          int       `json:"id"`
	EventID     int       `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
}

type WorkshopList struct {
	ID        int         `json:"id"`
	Title     string      `json:"title"`
	StartAt   time.Time   `json:"start_at"`
	EndAt     time.Time   `json:"end_at"`
	Workshops []*Workshop `json:"workshops"`
}

type WorkshopDetail struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	StartAt           time.Time `json:"start_at"`
	EndAt             time.Time `json:"end_at"`
	TotalReservations int       `json:"total_reservations"`
}
