package laundry

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Slot represents an available slot and corresponding machines
type Slot struct {
	Id       int       `db:"id"         json:"-"`
	Weekday  int       `db:"week_day"   json:"week_day"`
	Start    string    `db:"start_time" json:"start"`
	End      string    `db:"end_time"   json:"end"`
	Machines []Machine `                json:"machines"`
}

// SlotMachine represents the connection between one slot and one machine
type SlotMachine struct {
	Id        int `db:"id"          json:"-"`
	SlotId    int `db:"id_slots"    json:"id_slot"`
	MachineId int `db:"id_machines" json:"id_machine"`
}

// SlotWithBooker represents a slot and a possible booker for that slot
type SlotWithBooker struct {
	Slot
	Booker *Booker `json:"booker"`
}

// GetSlots will return a list of all slots and it's machines
func (l *Laundry) GetSlots() ([]Slot, error) {
	var slots []Slot

	slotSql := `SELECT * FROM slots`
	machineSql := ` SELECT m.* FROM machines AS m JOIN slots_machines AS sm
		ON sm.id_machines = m.id WHERE sm.id_slots = ?`

	rows, err := l.db.Queryx(slotSql)
	if err != nil {
		l.Logger.Errorf("Could not get slots: %s", err)
		return slots, ExtError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var s Slot
		if err := rows.StructScan(&s); err != nil {
			l.Logger.Errorf("Could not fetch row: %s", err)
			return slots, ExtError(err)
		}

		mRows, err := l.db.Queryx(machineSql, s.Id)
		if err != nil {
			l.Logger.Errorf("Could not get machines: %s", err)
			return slots, ExtError(err)
		}

		defer mRows.Close()

		for mRows.Next() {
			var m Machine
			if err := mRows.StructScan(&m); err != nil {
				l.Logger.Errorf("Could not fetch row: %s", err)
				return slots, ExtError(err)
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
func (l *Laundry) GetIntervalSchedule(start, end string) (map[time.Time][]SlotWithBooker, error) {
	sTime, eTime, err := dateIntervals(start, end)
	if err != nil {
		return nil, err
	}

	slots, _ := l.GetSlots()
	bookings, _ := l.GetBookingsInterval(start, end)

	var month = make(map[time.Time][]SlotWithBooker)

	for d := *sTime; d != (*eTime).AddDate(0, 0, 1); d = d.AddDate(0, 0, 1) {
		var fs []SlotWithBooker

		for _, s := range slots {
			var full SlotWithBooker

			// Check for slots this weekday
			if d.Weekday() == time.Weekday(s.Weekday) {
				full.Slot = s

				for _, b := range *bookings {
					// Check for bookings
					if b.BookDate == d && s.Start == b.Start {
						full.Booker = &b.Booker
					}
				}

				fs = append(fs, full)
			}
		}

		month[d] = fs
	}

	return month, nil
}
