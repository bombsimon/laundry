package laundry

import (
	"database/sql"
	"fmt"
	"time"

	//"github.com/davecgh/go-spew/spew"
	_ "github.com/jmoiron/sqlx"
)

type Booker struct {
	Id         int     `db:"id"         json:"-"`
	Identifier string  `db:"identifier" json:"identifier"` // Apartment number
	Name       *string `db:"name"       json:"name"`
	Email      *string `db:"email"      json:"email"`
	Phone      *string `db:"phone"      json:"phone"`
	Pin        *string `db:"pin"        json:"-"`
}

type Bookings struct {
	Id       int       `db:"id"        json:"-"`
	BookDate time.Time `db:"book_date" json:"book_date"`
	SlotId   int       `db:"id_slots"  json:"slot_id"`
	BookerId int       `db:"id_booker" json:"booker_id"`
}

type BookerBookingsRow struct {
	Booker
	Bookings
	Slot
	Machine
}

type BookerBookings struct {
	BookDate time.Time `json:"date"`
	Start    string    `json:"start"`
	End      string    `json:"end"`
	Booker   Booker    `json:"booker"`
	Machines []Machine `json:"machines"`
}

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

func (l *Laundry) GetBookings(b *Booker) (*[]BookerBookings, error) {
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

	defer rows.Close()

	var bookings = make(map[Bookings]BookerBookings)
	var machines = make(map[Bookings][]Machine)

	for rows.Next() {
		var bl = new(BookerBookingsRow)
		if err = rows.StructScan(bl); err != nil {
			l.Logger.Warnf("Could not gett bookings for booker with id %d: ", b.Id, err)
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

func (b *Booker) Notify() *Booker {
	// Send an email

	return b
}
