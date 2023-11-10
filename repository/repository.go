package repository

import (
	"database/sql"
	"time"

	"event-service/model"
)

type EventRepo struct {
	db *sql.DB
}

type EventRepository interface {
	GetActiveEventsWithPagination(limit, offset int) ([]*model.Event, error)
	GetActiveEventsCount() (int, error)
	GetEventByID(id int) (*model.EventDetail, error)
	GetWorkshopList(eventID int) (*model.WorkshopList, error)
	GetWorkshopDetail(workshopID int) (*model.WorkshopDetail, error)
	MakeReservation(workshopID int, name, email string) (*model.MakeReservationResponse, error)
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &EventRepo{db}
}

func (er *EventRepo) GetActiveEventsWithPagination(limit, offset int) ([]*model.Event, error) {
	query := `
        SELECT id, title, start_at, end_at
        FROM events
        WHERE start_at > NOW()
        ORDER BY start_at
        LIMIT ? OFFSET ?
    `

	rows, err := er.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var events []*model.Event
	for rows.Next() {
		event := &model.Event{}
		if err := rows.Scan(&event.ID, &event.Title, &event.StartAt, &event.EndAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (er *EventRepo) GetActiveEventsCount() (int, error) {
	query := `
		SELECT COUNT(*)
		FROM events
		WHERE start_at > NOW()
	`

	var count int
	if err := er.db.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (er *EventRepo) GetEventByID(id int) (*model.EventDetail, error) {
	query := `
		SELECT e.id, e.title, e.start_at, e.end_at, COUNT(w.id) AS total_workshops
		FROM events e
				 LEFT JOIN workshops w ON e.id = w.event_id
		WHERE e.id = ?
	`

	event := &model.EventDetail{}
	if err := er.db.
		QueryRow(query, id).
		Scan(
			&event.ID,
			&event.Title,
			&event.StartAt,
			&event.EndAt,
			&event.TotalWorkshop,
		); err != nil {
		return nil, err
	}

	return event, nil
}

func (er *EventRepo) GetWorkshopList(eventID int) (*model.WorkshopList, error) {
	query := `
		SELECT e.id AS event_id, 
		       e.title AS event_title, 
		       e.start_at AS event_start_at, 
		       e.end_at AS event_end_at,
		       w.id AS workshop_id, 
		       w.title AS workshop_title, 
		       w.description AS workshop_description,
		       w.start_at AS workshop_start_at, 
		       w.end_at AS workshop_end_at
		FROM events e
		JOIN workshops w ON e.id = w.event_id
		WHERE e.id = ? AND w.start_at > NOW()
	`

	rows, err := er.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	workshopList := &model.WorkshopList{}
	for rows.Next() {
		var (
			eventTitle                         string
			eventStartAt, eventEndAt           time.Time
			workshopID                         int
			workshopTitle, workshopDescription string
			workshopStartAt, workshopEndAt     time.Time
		)

		if err = rows.Scan(
			&eventID, &eventTitle, &eventStartAt, &eventEndAt,
			&workshopID, &workshopTitle, &workshopDescription,
			&workshopStartAt, &workshopEndAt,
		); err != nil {
			return nil, err
		}

		workshopList.ID = eventID
		workshopList.Title = eventTitle
		workshopList.StartAt = eventStartAt
		workshopList.EndAt = eventEndAt
		workshopList.Workshops = append(workshopList.Workshops, &model.Workshop{
			ID:          workshopID,
			Title:       workshopTitle,
			Description: workshopDescription,
			StartAt:     workshopStartAt,
			EndAt:       workshopEndAt,
		})
	}

	if len(workshopList.Workshops) == 0 {
		return nil, sql.ErrNoRows
	}

	return workshopList, nil
}

func (er *EventRepo) GetWorkshopDetail(workshopID int) (*model.WorkshopDetail, error) {
	query := `
		SELECT w.id,
		   w.title,
		   w.description,
		   w.start_at,
		   w.end_at,
		   COUNT(r.id) AS total_reservations
		FROM workshops w
				 LEFT JOIN reservations r ON w.id = r.workshop_id
		WHERE w.id = ?;
	`

	workshop := &model.WorkshopDetail{}
	if err := er.db.
		QueryRow(query, workshopID).
		Scan(
			&workshop.ID,
			&workshop.Title,
			&workshop.Description,
			&workshop.StartAt,
			&workshop.EndAt,
			&workshop.TotalReservations,
		); err != nil {
		return nil, err
	}

	return workshop, nil
}

func (er *EventRepo) insertReservation(workshopID int, name, email string) (*model.Reservation, error) {
	// Begin a transaction
	tx, err := er.db.Begin()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO reservations (name, email, workshop_id)
		VALUES (?, ?, ?)
	`

	res, err := tx.Exec(query, name, email, workshopID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Reservation{
		ID:    int(id),
		Name:  name,
		Email: email,
	}, nil
}

func (er *EventRepo) MakeReservation(workshopID int, name, email string) (*model.MakeReservationResponse, error) {
	reservation, err := er.insertReservation(workshopID, name, email)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT w.id AS workshop_id, w.title AS workshop_title, w.description AS workshop_description,
		       w.start_at AS workshop_start_at, w.end_at AS workshop_end_at,
		       e.id AS event_id, e.title AS event_title, e.start_at AS event_start_at, e.end_at AS event_end_at
		FROM workshops w
		JOIN events e ON w.event_id = e.id
		WHERE w.id = ?
	`

	var (
		workshopTitle, workshopDescription string
		workshopStartAt, workshopEndAt     time.Time
		eventID                            int
		eventTitle                         string
		eventStartAt, eventEndAt           time.Time
	)

	row := er.db.QueryRow(query, workshopID)
	if err = row.Scan(
		&workshopID,
		&workshopTitle,
		&workshopDescription,
		&workshopStartAt,
		&workshopEndAt,
		&eventID,
		&eventTitle,
		&eventStartAt,
		&eventEndAt,
	); err != nil {
		return nil, err
	}

	makeReservationResponse := &model.MakeReservationResponse{
		Event: &model.MakeReservationEvent{
			ID:      eventID,
			Title:   eventTitle,
			StartAt: eventStartAt,
			EndAt:   eventEndAt,
		},
		Workshop: &model.MakeReservationWorkshop{
			ID:          workshopID,
			Title:       workshopTitle,
			Description: workshopDescription,
			StartAt:     workshopStartAt,
			EndAt:       workshopEndAt,
		},
		Reservation: reservation,
	}

	return makeReservationResponse, nil
}
