package laundry

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Booker represents a booker
type Booker struct {
	Id         int     `db:"id"         json:"id"`
	Identifier string  `db:"identifier" json:"identifier"` // Apartment number
	Name       *string `db:"name"       json:"name"`
	Email      *string `db:"email"      json:"email"`
	Phone      *string `db:"phone"      json:"phone"`
	Pin        *string `db:"pin"        json:"-"`
}

// Bookings represents a booking
type Bookings struct {
	Id       int       `db:"id"        json:"-"`
	BookDate time.Time `db:"book_date" json:"book_date"`
	SlotId   int       `db:"id_slots"  json:"slot_id"`
	BookerId int       `db:"id_booker" json:"booker_id"`
}

// BookerBookingsRow represents the database struct to use when fetching
// bookings (slots, machines) and booker.
type BookerBookingsRow struct {
	Booker
	Bookings
	Slot
	Machine
}

// BookerBookings represents a booking including a Booker structure
type BookerBookings struct {
	BookDate time.Time `json:"date"`
	Start    string    `json:"start"`
	End      string    `json:"end"`
	Booker   Booker    `json:"booker"`
	Machines []Machine `json:"machines"`
}

// GetBooker will return a booker based on an id. If the booker is not found
// or an error fetching the booker occurs, an error will be returned.
func (l *Laundry) GetBooker(id int) (*Booker, error) {
	sqlStmt := `SELECT * FROM booker WHERE id = ?`
	stmt, err := l.db.Preparex(sqlStmt)

	defer stmt.Close()

	if err != nil {
		l.Logger.Errorf("Could not prepare statement: %s", err)
		return nil, ExtError(err)
	}

	var b Booker
	if err = stmt.QueryRowx(id).StructScan(&b); err == sql.ErrNoRows {
		l.Logger.Warnf("Booker with ID %d not found", id)
		return nil, NewError(fmt.Sprintf("Booker with id %d not found", id)).WithStatus(404)
	} else if err != nil {
		l.Logger.Errorf("Could not get row: %s", err)
		return nil, ExtError(err)
	}

	return &b, nil
}

// GetBookers will return a list of all bookers available
func (l *Laundry) GetBookers() ([]Booker, error) {
	sqlStmt := `SELECT * FROM booker`

	var bookers []Booker

	rows, err := l.db.Queryx(sqlStmt)
	if err != nil {
		l.Logger.Errorf("Could not get bookers")
		return bookers, ExtError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var b Booker
		if err := rows.StructScan(&b); err != nil {
			l.Logger.Errorf("Could not scan row: %s", err)
			return bookers, err
		}

		bookers = append(bookers, b)
	}

	return bookers, nil
}

// AddBooker will take a Booker strucutre and add it to the database.
// If the Booker is an existing Booker (or has an id), the id will be
// omitted and a possible copy of the booker will be created.
func (l *Laundry) AddBooker(b *Booker) (*Booker, error) {
	sqlStmt := `
	INSERT INTO booker (identifier, name, email, phone, pin)
	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := l.db.Preparex(sqlStmt)

	defer stmt.Close()

	if err != nil {
		l.Logger.Errorf("Could not create booker: %s", err)
		return nil, err
	}

	if _, err = stmt.Exec(b.Identifier, b.Name, b.Email, b.Phone, b.Pin); err != nil {
		l.Logger.Errorf("Could not create booker: %s", err)
		return nil, err
	}

	return b, nil
}

// UpdateBooker will take a Booker and update the row with corresponding
// id with the data in the Booker object.
func (l *Laundry) UpdateBooker(b *Booker) (*Booker, error) {
	sqlStmt := `
	UPDATE
		booker
	SET
		identifier = ?,
		name	   = ?,
		email      = ?,
		phone      = ?,
		pin        = ?
	WHERE
		id = ?
	`

	stmt, err := l.db.Preparex(sqlStmt)
	if err != nil {
		l.Logger.Warnf("Could not prepare statement: %s", err)
		return nil, err
	}

	if _, err = stmt.Exec(b.Identifier, b.Name, b.Email, b.Phone, b.Pin, b.Id); err != nil {
		l.Logger.Warnf("Could not update booker with ID %d: %s", b.Id, err)
		return nil, err
	}

	return b, nil
}

// RemoveBooker will take a Booker and remove the row with corresponding
// id in the database. A remove will cascade and remove belonging bookings
// and notifications.
func (l *Laundry) RemoveBooker(b *Booker) error {
	sqlStmt := `DELETE FROM booker WHERE id = ?`

	stmt, err := l.db.Preparex(sqlStmt)

	defer stmt.Close()

	if err != nil {
		l.Logger.Errorf("Could not prepare statement: %s", err)
		return NewError("Could not prepare statement").Add("Could not remove booker")
	}

	if _, err := stmt.Exec(b.Id); err != nil {
		l.Logger.Errorf("Could note remove booker")
		return ExtError(err).Add("Could not remove booker")
	}

	return nil
}

// ParseBookings will take an *sql.Rows and parse to a list of BookerBookings
func (l *Laundry) ParseBookings(rows *sqlx.Rows) (*[]BookerBookings, error) {
	defer rows.Close()

	var bookings = make(map[Bookings]BookerBookings)
	var machines = make(map[Bookings][]Machine)

	for rows.Next() {
		var bl = new(BookerBookingsRow)
		if err := rows.StructScan(bl); err != nil {
			l.Logger.Warnf("Could not gett bookings: ", err)
			return nil, ExtError(err)
		}

		booker := Booker{
			Identifier: bl.Booker.Identifier,
			Name:       bl.Booker.Name,
			Email:      bl.Booker.Email,
			Phone:      bl.Booker.Phone,
		}

		machine := Machine{
			Id:      bl.Machine.Id,
			Info:    bl.Machine.Info,
			Working: bl.Machine.Working,
		}

		br := BookerBookings{
			Booker:   booker,
			BookDate: bl.Bookings.BookDate,
			Start:    bl.Slot.Start,
			End:      bl.Slot.End,
		}

		bookings[bl.Bookings] = br
		machines[bl.Bookings] = append(machines[bl.Bookings], machine)
	}

	var retBookings []BookerBookings

	for k, v := range machines {
		bk, _ := bookings[k]
		bk.Machines = v
		bookings[k] = bk

		retBookings = append(retBookings, bookings[k])
	}

	return &retBookings, nil
}

// GetBookerBookings will take a Booker and return a set of BookerBookings
// for that Booker. The bookings returned will always be future bookings and
// not from the past.
func (l *Laundry) GetBookerBookings(b *Booker) (*[]BookerBookings, error) {
	sqlStmt := `
	SELECT
		b1.*, b2.*, s.*, m.*
	FROM
		bookings AS b1
	JOIN
		booker AS b2 ON
			b1.id_booker = b2.id
	JOIN
		slots AS s ON
			b1.id_slots = s.id
	JOIN
		slots_machines AS sm ON
			s.id = sm.id_slots
	JOIN
		machines AS m ON
			sm.id_machines = m.id
	WHERE
		b1.book_date >= DATE(NOW()) AND
		b2.id = ?
	`

	stmt, err := l.db.Preparex(sqlStmt)
	if err != nil {
		l.Logger.Warnf("Could not prepare statement: %s", err)
		return nil, ExtError(err)
	}

	rows, err := stmt.Queryx(b.Id)
	if err != nil {
		l.Logger.Errorf("Could not get bookings: %s", err)
		return nil, ExtError(err)
	}

	result, err := l.ParseBookings(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetBookingsInterval will return a list of BookerBookings between two given dates.
func (l *Laundry) GetBookingsInterval(start, end string) (*[]BookerBookings, error) {
	sTime, eTime, err := dateIntervals(start, end)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	SELECT
		b1.*, b2.*, s.*, m.*
	FROM
		bookings AS b1
	JOIN
		booker AS b2 ON
			b1.id_booker = b2.id
	JOIN
		slots AS s ON
			b1.id_slots = s.id
	JOIN
		slots_machines AS sm ON
			s.id = sm.id_slots
	JOIN
		machines AS m ON
			sm.id_machines = m.id
	WHERE
		b1.book_date >= DATE(?) AND
		b1.book_date <= DATE(?)
	`

	stmt, err := l.db.Preparex(sqlStmt)
	if err != nil {
		l.Logger.Warnf("Could not prepare statement: %s", err)
		return nil, ExtError(err)
	}

	rows, err := stmt.Queryx(sTime.String(), eTime.String())
	if err != nil {
		l.Logger.Errorf("Could not get bookings: %s", err)
		return nil, ExtError(err)
	}

	result, err := l.ParseBookings(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Notify will send a notification to the Booker
func (b *Booker) Notify() *Booker {
	// Send an email

	return b
}
