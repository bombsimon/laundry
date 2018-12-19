package pkg

import "time"

// Slot represents an available slot and corresponding machines
type Slot struct {
	ID       int        `db:"id"         json:"id"`
	Weekday  int        `db:"week_day"   json:"week_day"`
	Start    string     `db:"start_time" json:"start"`
	End      string     `db:"end_time"   json:"end"`
	Machines []*Machine `db:"-"          json:"machines"`
}

// SlotWithBooker represents a slot and a possible booker for that slot
type SlotWithBooker struct {
	Slot
	Booker *Booker `json:"booker"`
}

// Schedule is a list of all slots mapped by time to a list of bookings.
type Schedule map[time.Time][]*SlotWithBooker

// SlotHandler is the interface to implement to handle slots.
type SlotHandler interface {
	GetSlot(id int) (*Slot, error)
	GetSlots() ([]*Slot, error)
	AddSlot(slot *Slot) (*Slot, error)
	UpdateSlot(slotID int, toUpdate *Slot) (*Slot, error)
	RemoveSlot(slot *Slot) error
	RemoveSlotByID(id int) error
	GetSchedule(start, end time.Time) (Schedule, error)
}
