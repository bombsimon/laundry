package laundry

import (
	_ "github.com/go-sql-driver/mysql"
)

type Slot struct {
	Id       int       `db:"id"         json:"-"`
	WeekDay  int       `db:"week_day"   json:"week_day"`
	Start    string    `db:"start_time" json:"start"`
	End      string    `db:"end_time"   json:"end"`
	Machines []Machine `                json:"machines"`
	Booker   *Booker   `                json:"booker"`
	Notify   []*Booker `                json:"notify"`
}

func (l *Laundry) Release(s *Slot) *Slot {
	// Notify watchers when slot is released
	for _, booker := range s.Notify {
		booker.Notify()
	}

	sql := "UPDATE slots SET id_booker = NULL WHERE id = ?"
	_, err := l.db.Queryx(sql, s.Id)
	if err != nil {
		// Handle update problem
	}

	// Slot no longer has a booker
	s.Booker = nil
	return s
}
