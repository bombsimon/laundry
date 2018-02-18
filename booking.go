package laundry

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bombsimon/laundry/database"
	"github.com/bombsimon/laundry/errors"
	"github.com/jmoiron/sqlx"
)

// BookingsSeach represents searchable booking parameters
type BookingsSearch struct {
	Start  time.Time
	End    time.Time
	Booker *Booker
}

// Booker represents a booker
type Booker struct {
	Id         int        `db:"id"         json:"id"`
	Identifier string     `db:"identifier" json:"identifier"` // Apartment number
	Name       NullString `db:"name"       json:"name"`
	Email      NullString `db:"email"      json:"email"`
	Phone      NullString `db:"phone"      json:"phone"`
	Pin        NullString `db:"pin"        json:"-"`
}

// Bookings represents a booking
type Bookings struct {
	Id       int       `db:"id"        json:"id"`
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
	Slot     Slot      `json:"slot"`
	Booker   Booker    `json:"booker"`
	Machines []Machine `json:"machines"`
}

// GetBooker will return a booker based on an id. If the booker is not found
// or an error fetching the booker occurs, an error will be returned.
func GetBooker(id int) (*Booker, *errors.LaundryError) {
	db := database.GetConnection()

	var b Booker
	if err := db.QueryRowx("SELECT * FROM booker WHERE id = ?", id).StructScan(&b); err == sql.ErrNoRows {
		return nil, errors.New("Booker with id %d not found", id).WithStatus(http.StatusNotFound)
	} else if err != nil {
		return nil, errors.New("Could not get row").CausedBy(err)
	}

	return &b, nil
}

// GetBookers will return a list of all bookers available
func GetBookers() ([]Booker, *errors.LaundryError) {
	db := database.GetConnection()

	var bookers []Booker
	rows, err := db.Queryx("SELECT * FROM booker")
	if err != nil {
		return bookers, errors.New("Could not get bookers").CausedBy(err)
	}

	defer rows.Close()

	for rows.Next() {
		var b Booker
		if err := rows.StructScan(&b); err != nil {
			return bookers, errors.New("Could not scan row").CausedBy(err)
		}

		bookers = append(bookers, b)
	}

	return bookers, nil
}

// AddBooker will take a Booker strucutre and add it to the database.
// If the Booker is an existing Booker (or has an id), the id will be
// omitted and a possible copy of the booker will be created.
func AddBooker(b *Booker) (*Booker, *errors.LaundryError) {
	if b.Identifier == "" {
		return nil, errors.New("Missing identifier in request").WithStatus(http.StatusBadRequest)
	}

	db := database.GetConnection()

	query := "INSERT INTO booker (identifier, name, email, phone, pin) VALUES (?, ?, ?, ?, ?)"

	row, err := db.Exec(query, b.Identifier, b.Name, b.Email, b.Phone, b.Pin)
	if err != nil {
		return nil, errors.New("Could not create booker").CausedBy(err)
	}

	lastId, err := row.LastInsertId()
	if err != nil {
		return nil, errors.New(err)
	}

	b.Id = int(lastId)

	return b, nil
}

// UpdateBooker will take a Booker and update the row with corresponding
// id with the data in the Booker object.
func UpdateBooker(bookerId int, ub *Booker) (*Booker, *errors.LaundryError) {
	b, berr := GetBooker(bookerId)
	if berr != nil {
		return nil, berr
	}

	b.Phone = ub.Phone
	b.Email = ub.Email

	db := database.GetConnection()

	query := `
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

	if _, err := db.Exec(query, b.Identifier, b.Name, b.Email, b.Phone, b.Pin, b.Id); err != nil {
		return nil, errors.New("Could not update booker with ID %d", b.Id).CausedBy(err)
	}

	return b, nil
}

// RemoveBooker will take a Booker and remove the row with corresponding
// id in the database. A remove will cascade and remove belonging bookings
// and notifications.
func RemoveBooker(b *Booker) *errors.LaundryError {
	db := database.GetConnection()

	if _, err := db.Exec("DELETE FROM booker WHERE id = ?", b.Id); err != nil {
		return errors.New("Could not remove booker").CausedBy(err)
	}

	return nil
}

// RemoveBookerById will remove a booker by the booker id
func RemoveBookerById(id int) *errors.LaundryError {
	b, err := GetBooker(id)
	if err != nil {
		return err
	}

	return RemoveBooker(b)
}

// GetBookerBookings will take a Booker and return a set of BookerBookings
// for that Booker. The bookings returned will always be future bookings and
// not from the past.
func GetBookerBookings(b *Booker) (*[]BookerBookings, *errors.LaundryError) {
	bs := BookingsSearch{
		time.Now(),
		time.Now().AddDate(10, 0, 0),
		b,
	}

	return SearchBookings(bs)
}

// GetBookerBookingsById will take a booker id and return bookings if the id is
// bound to a booker
func GetBookerBookingsById(id int) (*[]BookerBookings, *errors.LaundryError) {
	b, err := GetBooker(id)
	if err != nil {
		return nil, err
	}

	return GetBookerBookings(b)
}

// SearchBookings will return a list of BookerBookings based on passed search criteria
func SearchBookings(bs BookingsSearch) (*[]BookerBookings, *errors.LaundryError) {
	db := database.GetConnection()

	var (
		query      string
		conditions []interface{}
	)

	query = `
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

	conditions = append(conditions, bs.Start.String(), bs.End.String())

	if bs.Booker != nil {
		query += `AND b2.id = ?`

		conditions = append(conditions, bs.Booker.Id)
	}

	rows, err := db.Queryx(query, conditions...)
	if err != nil {
		return nil, errors.New("Could not get bookings").CausedBy(err)
	}

	result, parseErr := parseBookings(rows)
	if parseErr != nil {
		return nil, errors.New(parseErr)
	}

	return result, nil
}

// parseBookings will take an *sql.Rows and parse to a list of BookerBookings
func parseBookings(rows *sqlx.Rows) (*[]BookerBookings, *errors.LaundryError) {
	defer rows.Close()

	var bookings = make(map[Bookings]BookerBookings)
	var machines = make(map[Bookings][]Machine)

	for rows.Next() {
		var bl = new(BookerBookingsRow)
		if err := rows.StructScan(bl); err != nil {
			return nil, errors.New("Could not get bookings").CausedBy(err)
		}

		booker := Booker{
			Id:         bl.Booker.Id,
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

		slot := Slot{
			Id:      bl.Slot.Id,
			Weekday: bl.Slot.Weekday,
			Start:   bl.Slot.Start,
			End:     bl.Slot.End,
		}

		br := BookerBookings{
			BookDate: bl.Bookings.BookDate,
			Booker:   booker,
			Slot:     slot,
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
