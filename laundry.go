/*
Laundry is a service used to manage laundry bookings, primarily in realestates
with a shared laundry room.

The service includes interfaces for every part in the system and has an beloning
RESTful API to use in combination with a GUI or front end service.
*/
package laundry

import (
	"time"

	"github.com/bombsimon/laundry/errors"
	_ "github.com/go-sql-driver/mysql"
)

// dateIntervals is a generic function that takes two strings, parses them
// as time.Time objects and makes sure the start time does not occure after
// the end time.
func dateIntervals(start, end string) (*time.Time, *time.Time, error) {
	sTime, _ := time.Parse("2006-01-02", start)
	eTime, _ := time.Parse("2006-01-02", end)

	if sTime.After(eTime) {
		return nil, nil, errors.New("Start time cannot be after end time")
	}

	return &sTime, &eTime, nil
}
