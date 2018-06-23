package laundry

import (
	"net/http"
	"time"

	"github.com/bombsimon/laundry/database"
	"github.com/bombsimon/laundry/errors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/doug-martin/goqu.v4"
	// MySQL driver for goqu
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
)

// BookingsSearch represents searchable booking parameters
type BookingsSearch struct {
	Start  time.Time
	End    time.Time
	Booker *Booker
}

// Booker represents a booker
type Booker struct {
	ID         int        `db:"id"         json:"id"`
	Identifier string     `db:"identifier" json:"identifier"` // Apartment number
	Name       NullString `db:"name"       json:"name"`
	Email      NullString `db:"email"      json:"email"`
	Phone      NullString `db:"phone"      json:"phone"`
	Pin        NullString `db:"pin"        json:"-"`
}

// Bookings represents a booking
type Bookings struct {
	ID       int       `db:"id"        json:"id"`
	BookDate time.Time `db:"book_date" json:"book_date"`
	SlotID   int       `db:"id_slots"  json:"slot_id"`
	BookerID int       `db:"id_booker" json:"booker_id"`
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
	db := database.GetGoqu()

	var b Booker
	found, err := db.From("booker").Where(goqu.Ex{
		"id": id,
	}).ScanStruct(&b)

	if err != nil {
		return nil, errors.New("Could not get row").CausedBy(err)
	}

	if !found {
		return nil, errors.New("Booker with id %d not found", id).WithStatus(http.StatusNotFound)
	}

	return &b, nil
}

// GetBookers will return a list of all bookers available
func GetBookers() ([]Booker, *errors.LaundryError) {
	db := database.GetGoqu()

	var bookers []Booker
	if err := db.From("booker").ScanStructs(&bookers); err != nil {
		return bookers, errors.New("Could not get bookers").CausedBy(err)
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

	db := database.GetGoqu()

	insert := db.From("booker").Insert(goqu.Record{
		"identifier": b.Identifier,
		"name":       b.Name,
		"email":      b.Email,
		"phone":      b.Phone,
		"pin":        b.Pin,
	})

	row, err := insert.Exec()
	if err != nil {
		return nil, errors.New("Could not create booker").CausedBy(err)
	}

	lastID, err := row.LastInsertId()
	if err != nil {
		return nil, errors.New(err)
	}

	b.ID = int(lastID)

	return b, nil
}

// UpdateBooker will take a Booker and update the row with corresponding
// id with the data in the Booker object.
func UpdateBooker(bookerID int, ub *Booker) (*Booker, *errors.LaundryError) {
	b, berr := GetBooker(bookerID)
	if berr != nil {
		return nil, berr
	}

	b.Name = ub.Name
	b.Email = ub.Email
	b.Phone = ub.Phone

	db := database.GetGoqu()

	update := db.From("booker").
		Where(goqu.Ex{
			"id": b.ID,
		}).
		Update(goqu.Record{
			"identifier": b.Identifier,
			"name":       b.Name,
			"email":      b.Email,
			"phone":      b.Phone,
			"pin":        b.Pin,
		})

	if _, err := update.Exec(); err != nil {
		return nil, errors.New("Could not update booker with ID %d", b.ID).CausedBy(err)
	}

	return b, nil
}

// RemoveBooker will take a Booker and remove the row with corresponding
// id in the database. A remove will cascade and remove belonging bookings
// and notifications.
func RemoveBooker(b *Booker) *errors.LaundryError {
	db := database.GetGoqu()

	delete := db.From("booker").
		Where(goqu.Ex{
			"id": b.ID,
		}).
		Delete()

	if _, err := delete.Exec(); err != nil {
		return errors.New("Could not remove booker").CausedBy(err)
	}

	return nil
}

// RemoveBookerByID will remove a booker by the booker id
func RemoveBookerByID(id int) *errors.LaundryError {
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

// GetBookerBookingsByID will take a booker id and return bookings if the id is
// bound to a booker
func GetBookerBookingsByID(id int) (*[]BookerBookings, *errors.LaundryError) {
	b, err := GetBooker(id)
	if err != nil {
		return nil, err
	}

	return GetBookerBookings(b)
}

// SearchBookings will return a list of BookerBookings based on passed search criteria
// TODO: This should not be inflated like this, the query is bad.
func SearchBookings(bs BookingsSearch) (*[]BookerBookings, *errors.LaundryError) {
	db := database.GetGoqu()

	query := db.From("bookings").
		Select("bookings.*", "booker.*", "slots.*", "machines.*").
		LeftJoin(goqu.I("booker"), goqu.On(goqu.I("booker.id").Eq(goqu.I("bookings.id_booker")))).
		LeftJoin(goqu.I("slots"), goqu.On(goqu.I("slots.id").Eq(goqu.I("bookings.id_slots")))).
		LeftJoin(goqu.I("slots_machines"), goqu.On(goqu.I("slots_machines.id_slots").Eq(goqu.I("slots.id")))).
		LeftJoin(goqu.I("machines"), goqu.On(goqu.I("machines.id").Eq(goqu.I("slots_machines.id_machines")))).
		Where(
			goqu.L("bookings.book_date >= DATE(?)", bs.Start.String()),
			goqu.L("bookings.book_date <= DATE(?)", bs.End.String()),
		).
		Prepared(true)

	if bs.Booker != nil {
		query = query.
			Where(
				goqu.I("booker.id").Eq(bs.Booker.ID),
			).
			Prepared(true)
	}

	sql, args, _ := query.ToSql()

	sqlxDb := database.GetConnection()
	rows, err := sqlxDb.Queryx(sql, args...)
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
			ID:         bl.Booker.ID,
			Identifier: bl.Booker.Identifier,
			Name:       bl.Booker.Name,
			Email:      bl.Booker.Email,
			Phone:      bl.Booker.Phone,
		}

		machine := Machine{
			ID:      bl.Machine.ID,
			Info:    bl.Machine.Info,
			Working: bl.Machine.Working,
		}

		slot := Slot{
			ID:      bl.Slot.ID,
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
