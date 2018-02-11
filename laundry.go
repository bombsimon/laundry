/*
Laundry is a service used to manage laundry bookings, primarily in realestates
with a shared laundry room.

The service includes interfaces for every part in the system and has an beloning
RESTful API to use in combination with a GUI or front end service.

	laundry := laundry.New("/path/to/service-config.yaml")
	bookers, err := laundry.GetBookers()
	if err != nil {
		laundry.Logger.Warnf("Something is not right: %s", err)

		// Appropriate HTTP status - err.(*laundry.LaundryError).Status()
		// Marshal error to JSON   - err.(*laundry.LaundryError).AsJSON()
	}
*/
package laundry

import (
	"os"
	"time"

	"github.com/bombsimon/laundry/config"
	"github.com/bombsimon/laundry/database"
	"github.com/bombsimon/laundry/errors"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// Laundry represents a laundry service with a database handler, a logger
// and configuration.
type Laundry struct {
	Logger *log.Entry
	Config *config.Configuration
}

// New will take a string with a path to the configuration file and setup a
// new Laundry object. If the configuration file does not exist an object will
// still be returned with empty configuration.
// There are two environment variables to use to override the configuration file:
//  LAUNDRY_DSN         - DSN to use when connecting to database
//  LAUNDRY_HTTP_LISTEN - The host/port to listen on when running the API
func New(configFile string) *Laundry {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	logger := log.WithFields(log.Fields{})

	conf, err := config.New(configFile)
	if err != nil {
		logger.Warnf("Configuration file missing - trying to proceed with unknown result")
		conf = &config.Configuration{}
	}

	database.SetupConnection(conf.Database)

	if os.Getenv("LAUNDRY_HTTP_LISTEN") != "" {
		conf.HTTP.Listen = os.Getenv("LAUNDRY_HTTP_LISTEN")
	}

	l := Laundry{
		Logger: logger,
		Config: conf,
	}

	return &l
}

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
