package model

import (
	"time"
)

type Reservation struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MakeReservationWorkshop struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
}

type MakeReservationEvent struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

type MakeReservationResponse struct {
	Event       *MakeReservationEvent    `json:"event"`
	Workshop    *MakeReservationWorkshop `json:"workshop"`
	Reservation *Reservation             `json:"reservation"`
}
