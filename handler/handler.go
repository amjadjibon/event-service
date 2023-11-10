package handler

import (
	"event-service/repository"
)

type EventHandler struct {
	repo repository.EventRepository
}

func NewEventHandler(repo repository.EventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}
