/*
Laundry is a service used to manage laundry bookings, primarily in realestates
with a shared laundry room.

The service includes interfaces for every part in the system and has an beloning
RESTful API to use in combination with a GUI or front end service.
*/
package laundry

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/bombsimon/laundry/errors"
	_ "github.com/go-sql-driver/mysql"
)

// NullString represents an embedded sql.NullString on which we
// can implement a custom JSON marshaller
type NullString struct {
	sql.NullString
}

// MarshalJSON will make sure NullStrings are marshalled correct
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}

	return []byte("null"), nil
}

func dateIntervals(start, end string) (*time.Time, *time.Time, *errors.LaundryError) {
	return interval("2006-01-02", start, end, true)
}

func timeIntervals(start, end string) (*time.Time, *time.Time, *errors.LaundryError) {
	return interval("15:04:05", start, end, true)
}

func interval(format, start, end string, ordered bool) (*time.Time, *time.Time, *errors.LaundryError) {
	sTime, err := time.Parse(format, start)
	if err != nil {
		return nil, nil, errors.New("Invalid start time").CausedBy(err)
	}

	eTime, err := time.Parse(format, end)
	if err != nil {
		return nil, nil, errors.New("Invalid end time").CausedBy(err)
	}

	// Check that start- and end time was sent in order
	if ordered {
		if sTime.After(eTime) {
			return nil, nil, errors.New("Start time cannot be after end time")
		}
	}

	return &sTime, &eTime, nil
}
