package models

import (
	"ems_backend_go/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	CreatedBy   int
}

var events = []Event{}

func (e Event) Save() error {
	query := `INSERT INTO 
	events(name,description,location,dateTime,user_id) 
	VALUES(?,?,?,?,?)`
	stmnt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmnt.Close()
	res, err := stmnt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.CreatedBy)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.CreatedBy)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func DeleteAllEvents() error {
	query := `DELETE FROM events`
	_, err := db.DB.Exec(query)
	if err != nil {
		return err
	}
	query2 := `DELETE FROM sqlite_sequence WHERE name='events'`
	_, err = db.DB.Exec(query2)
	if err != nil {
		return err
	}
	return nil
}
