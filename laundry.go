package laundry

import (
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Laundry struct {
	db     *sqlx.DB
	Logger *log.Entry
}

func New(dsn string) *Laundry {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	logger := log.WithFields(log.Fields{})

	db := mysqlConnect(dsn, 5, logger)

	l := Laundry{
		Logger: logger,
		db:     db,
	}

	return &l
}

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

func dateIntervals(start, end string) (*time.Time, *time.Time, error) {
	sTime, _ := time.Parse("2006-01-02", start)
	eTime, _ := time.Parse("2006-01-02", end)

	if sTime.After(eTime) {
		return nil, nil, NewError("Start time cannot be after end time")
	}

	return &sTime, &eTime, nil
}
