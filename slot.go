package laundry

import (
	"net/http"
	"time"

	"github.com/bombsimon/laundry/database"
	"github.com/bombsimon/laundry/errors"
	"github.com/bombsimon/laundry/log"
	// MySQL driver for goqu
	goqu "gopkg.in/doug-martin/goqu.v4"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
)

// Slot represents an available slot and corresponding machines
type Slot struct {
	ID       int       `db:"id"         json:"id"`
	Weekday  int       `db:"week_day"   json:"week_day"`
	Start    string    `db:"start_time" json:"start"`
	End      string    `db:"end_time"   json:"end"`
	Machines []Machine `db:"-"          json:"machines"`
}

// SlotWithBooker represents a slot and a possible booker for that slot
type SlotWithBooker struct {
	Slot
	Booker *Booker `json:"booker"`
}

// GetSlots will return a list of all slots and it's machines
func GetSlots() ([]Slot, *errors.LaundryError) {
	db := database.GetGoqu()
	var slots []Slot

	if err := db.From("slots").ScanStructs(&slots); err != nil {
		return slots, errors.New("Could not get slots").CausedBy(err)
	}

	for i, slot := range slots {
		var machines []Machine

		err := db.From("machines").
			Select("machines.*").
			LeftJoin(goqu.I("slots_machines"), goqu.On(goqu.I("slots_machines.id_machines").Eq(goqu.I("machines.id")))).
			Where(
				goqu.I("slots_machines.id_slots").Eq(slot.ID),
			).ScanStructs(&machines)

		if err != nil {
			log.GetLogger().Errorf("Could not get machines: %s", err)
			return slots, errors.New(err)
		}

		slots[i].Machines = machines

	}

	return slots, nil
}

// AddSlot will create a new slot
func AddSlot(s *Slot) (*Slot, *errors.LaundryError) {
	if err := validSlot(s); err != nil {
		return nil, err
	}

	db := database.GetGoqu()

	insert := db.From("slots").Insert(goqu.Record{
		"week_day":   s.Weekday,
		"start_time": s.Start,
		"end_time":   s.End,
	})

	row, err := insert.Exec()
	if err != nil {
		return nil, errors.New("Could not create slot").CausedBy(err)
	}

	lastID, err := row.LastInsertId()
	if err != nil {
		return nil, errors.New(err)
	}

	s.ID = int(lastID)

	return s, nil
}

// GetSlot will return one slot
func GetSlot(slotID int) (*Slot, *errors.LaundryError) {
	db := database.GetGoqu()

	var s Slot
	found, err := db.From("slots").Where(goqu.Ex{
		"id": slotID,
	}).ScanStruct(&s)

	if err != nil {
		return nil, errors.New("Could not get row").CausedBy(err)
	}

	if !found {
		return nil, errors.New("Slot with id %d not found", slotID).WithStatus(http.StatusNotFound)
	}

	return &s, nil
}

// UpdateSlot will update an existing slot
func UpdateSlot(slotID int, s *Slot) (*Slot, *errors.LaundryError) {
	if err := validSlot(s); err != nil {
		return nil, err
	}

	slot, err := GetSlot(slotID)
	if err != nil {
		return nil, err
	}

	slot.Weekday = s.Weekday
	slot.Start = s.Start
	slot.End = s.End

	db := database.GetGoqu()

	update := db.From("slots").Where(goqu.Ex{
		"id": slot.ID,
	}).
		Update(goqu.Record{
			"week_day":   slot.Weekday,
			"start_time": slot.Start,
			"end_time":   slot.End,
		})

	if _, err := update.Exec(); err != nil {
		return nil, errors.New("Could not update slot with id %d", slot.ID).CausedBy(err)
	}

	return slot, nil
}

// RemoveSlot will remove an existing slot
func RemoveSlot(s *Slot) *errors.LaundryError {
	db := database.GetGoqu()

	delete := db.From("slots").Where(goqu.Ex{
		"id": s.ID,
	}).Delete()

	if _, err := delete.Exec(); err != nil {
		return errors.New("Could not remove slot with id %d", s.ID).CausedBy(err)
	}

	return nil
}

// RemoveSlotByID will remove a slot by a aslot id
func RemoveSlotByID(id int) *errors.LaundryError {
	slot, err := GetSlot(id)
	if err != nil {
		return err
	}

	return RemoveSlot(slot)
}

func validSlot(s *Slot) *errors.LaundryError {
	// Valid day provided
	switch s.Weekday {
	case 0, 1, 2, 3, 4, 5, 6:
		// Valid day
	default:
		return errors.New("Invalid weekday").WithStatus(http.StatusBadRequest)
	}

	// Valid start- and end time provided
	if _, _, err := timeIntervals(s.Start, s.End); err != nil {
		return err
	}

	return nil
}

// GetIntervalSchedule will return a schedule between a given start- and end time.
// A map for each day will be returned holding a list of slots and possible bookers
// for the given slot.
func GetIntervalSchedule(start, end string) (map[time.Time][]SlotWithBooker, *errors.LaundryError) {
	sTime, eTime, err := dateIntervals(start, end)
	if err != nil {
		return nil, errors.New(err)
	}

	// All slots in the system
	slots, _ := GetSlots()

	// All bookings in the system
	bookings, sErr := SearchBookings(BookingsSearch{*sTime, *eTime, nil})
	if sErr != nil {
		return nil, sErr
	}

	var month = make(map[time.Time][]SlotWithBooker)

	// Iterate from start date, add one day each iteration until we're at the end date
	for d := *sTime; d != (*eTime).AddDate(0, 0, 1); d = d.AddDate(0, 0, 1) {
		// SlotWithBooker is a slot in any given date which also includes a booker (if booker)
		// Each day may have multiple slots with one booker each
		var fs []SlotWithBooker

		// Iterate over all slots for every given date
		for _, s := range slots {
			// Each slot is bound to a week day. If the slot isn't on the same week day
			// as the current iteration, ignore it
			if d.Weekday() != time.Weekday(s.Weekday) {
				continue
			}

			// Add the current slot to the date we're at in our iterator

			var full = SlotWithBooker{
				Slot: s,
			}

			// Iterate over all bookings and see if any of them are at this current day
			// with the same start time as the slot
			// TODO: This is crap and high complexity - fix
			for _, b := range *bookings {
				// If the booking is on the same date as the iterator and the booking start time
				// is the same as the slot - add it to the result
				if b.BookDate == d && s.Start == b.Slot.Start {
					full.Booker = &b.Booker
				}
			}

			fs = append(fs, full)
		}

		// Add all found slots, with or without booker, to the current date
		month[d] = fs
	}

	return month, nil
}
