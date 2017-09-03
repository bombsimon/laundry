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
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// Laundry represents a laundry service with a database handler, a logger
// and configuration.
type Laundry struct {
	db     *sqlx.DB
	Logger *log.Entry
	Config *Configuration
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

	config, err := NewConfig(configFile)
	if err != nil {
		logger.Warnf("Configuration file missing - trying to proceed with unknown result")
		config = &Configuration{}
	}

	dsn := os.Getenv("LAUNDRY_DSN")
	if dsn == "" {
		db := config.Database
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=1", db.Username, db.Password, db.Host, db.Port, db.Database)
	}

	if os.Getenv("LAUNDRY_HTTP_LISTEN") != "" {
		config.HTTP.Listen = os.Getenv("LAUNDRY_HTTP_LISTEN")
	}

	db := mysqlConnect(dsn, 5, logger)

	l := Laundry{
		Logger: logger,
		db:     db,
		Config: config,
	}

	return &l
}

// mysqlConnect will take a DSN and a retry count and try to connect to the
// database that many times. This menas that it will take a long time
// to return a Laundry object if it's not possible to connect to the database.
func mysqlConnect(dsn string, retries int, logger *log.Entry) *sqlx.DB {
	for i := retries; i >= 0; i-- {
		db, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			logger.Infof("Could not connect. Retrying in %d seconds. Reason: %s", 5, err)
			time.Sleep(time.Second * 5)
			continue
		}

		return db
	}

	return nil
}

// dateIntervals is a generic function that takes two strings, parses them
// as time.Time objects and makes sure the start time does not occure after
// the end time.
func dateIntervals(start, end string) (*time.Time, *time.Time, error) {
	sTime, _ := time.Parse("2006-01-02", start)
	eTime, _ := time.Parse("2006-01-02", end)

	if sTime.After(eTime) {
		return nil, nil, NewError("Start time cannot be after end time")
	}

	return &sTime, &eTime, nil
}
