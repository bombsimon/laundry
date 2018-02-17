package laundry

import (
	"time"

	"github.com/bombsimon/laundry/database"
	"github.com/bombsimon/laundry/errors"
	"github.com/bombsimon/laundry/log"
	_ "github.com/go-sql-driver/mysql"
)

// Slot represents an available slot and corresponding machines
type Slot struct {
	Id       int       `db:"id"         json:"id"`
	Weekday  int       `db:"week_day"   json:"week_day"`
	Start    string    `db:"start_time" json:"start"`
	End      string    `db:"end_time"   json:"end"`
	Machines []Machine `                json:"machines"`
}

// SlotWithBooker represents a slot and a possible booker for that slot
type SlotWithBooker struct {
	Slot
	Booker *Booker `json:"booker"`
}

// GetSlots will return a list of all slots and it's machines
func GetSlots() ([]Slot, *errors.LaundryError) {
	var slots []Slot

	db := database.GetConnection()
	slotSql := `SELECT * FROM slots`
	machineSql := ` SELECT m.* FROM machines AS m JOIN slots_machines AS sm
		ON sm.id_machines = m.id WHERE sm.id_slots = ?`

	rows, err := db.Queryx(slotSql)
	if err != nil {
		return slots, errors.New("Could not get slots").CausedBy(err)
	}

	defer rows.Close()

	for rows.Next() {
		var s Slot
		if err := rows.StructScan(&s); err != nil {
			log.GetLogger().Errorf("Could not fetch row: %s", err)
			return slots, errors.New(err)
		}

		// TODO: Don't query the database for each slot - this should be a JOIN
		mRows, err := db.Queryx(machineSql, s.Id)
		if err != nil {
			log.GetLogger().Errorf("Could not get machines: %s", err)
			return slots, errors.New(err)
		}

		defer mRows.Close()

		for mRows.Next() {
			var m Machine
			if err := mRows.StructScan(&m); err != nil {
				log.GetLogger().Errorf("Could not fetch row: %s", err)
				return slots, errors.New(err)
			}

			s.Machines = append(s.Machines, m)
		}

		slots = append(slots, s)
	}

	return slots, nil
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
	bookings, _ := SearchBookings(BookingsSearch{*sTime, *eTime, nil})

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
				if b.BookDate == d && s.Start == b.Start {
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
