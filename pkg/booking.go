package pkg

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

// Booker represents a booker
type Booker struct {
	ID         int         `db:"id"         json:"id"`
	Identifier string      `db:"identifier" json:"identifier"` // Apartment number
	Name       null.String `db:"name"       json:"name"`
	Email      null.String `db:"email"      json:"email"`
	Phone      null.String `db:"phone"      json:"phone"`
	Pin        null.String `db:"pin"        json:"-"`
}

// Bookings represents a booking
type Bookings struct {
	ID       int       `db:"id"        json:"id"`
	BookDate time.Time `db:"book_date" json:"book_date"`
	SlotID   int       `db:"id_slots"  json:"slot_id"`
	BookerID int       `db:"id_booker" json:"booker_id"`
}

// BookerBookings represents a booking including a Booker structure
type BookerBookings struct {
	BookDate time.Time `json:"date"`
	Slot     Slot      `json:"slot"`
	Booker   Booker    `json:"booker"`
	Machines []Machine `json:"machines"`
}

// BookingFindCriteria is the criteria that can be present to search for a
// booking.
type BookingFindCriteria struct {
	BookerID int       `json:"booker_id"`
	SlotID   int       `json:"slot_id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

// BookingHandler is the interface that implements all feature to book a slot.
type BookingHandler interface {
	AddBooker(booker *Booker) (*Booker, error)
	FindBooking(fb BookingFindCriteria) ([]*BookerBookings, error)
	GetBooker(id int) (*Booker, error)
	GetBookerBookings(b *Booker) ([]*BookerBookings, error)
	GetBookerBookingsByID(id int) ([]*BookerBookings, error)
	GetBookers() ([]*Booker, error)
	RemoveBooker(b *Booker) error
	RemoveBookerByID(bookerID int) error
	UpdateBooker(bookerID int, toUpdate *Booker) (*Booker, error)
}
