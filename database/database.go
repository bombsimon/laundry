package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bombsimon/laundry/config"
	"github.com/bombsimon/laundry/log"
	"github.com/jmoiron/sqlx"
	goqu "gopkg.in/doug-martin/goqu.v4"
)

var (
	connectionPool []dbConnection
)

type dbConnection struct {
	db             *sqlx.DB
	simpleDb       *sql.DB
	reconnected    bool
	reconnectCount int
}

// SetupConnection will setup connections and store them in a
// connection pool
func SetupConnection(config config.Database) {
	logger := log.GetLogger()

	dsn := os.Getenv("LAUNDRY_DSN")
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=1", config.Username, config.Password, config.Host, config.Port, config.Database)
	}

	for i := 1; i <= config.PoolSize; i++ {
		var dbConn dbConnection

		for j := config.RetryCount; j >= 0; j-- {
			logger.Info("connection to databas")

			db, err := sql.Open("mysql", dsn)

			if err != nil {
				logger.Warnf("could not connect att attempt %d: %s", j, err)
				time.Sleep(time.Second * 5)
				continue
			}

			dbConn.simpleDb = db
			dbConn.db = sqlx.NewDb(db, "mysql")
		}

		go monitorConnection(dbConn)

		connectionPool = append(connectionPool, dbConn)
	}
}

// GetConnection will return a connection to the database
func GetConnection() *sqlx.DB {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	return connectionPool[r.Intn(len(connectionPool))].db
}

// GetSimpleConnection will return a simple SQL connection.
func GetSimpleConnection() *sql.DB {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	return connectionPool[r.Intn(len(connectionPool))].simpleDb
}

// GetGoqu will return a goqu type for goqu queries.
func GetGoqu() *goqu.Database {
	return goqu.New("mysql", GetSimpleConnection())
}

func monitorConnection(d dbConnection) {
	for {
		if err := d.db.Ping(); err != nil {
			// Should reconnect
			d.reconnected = true
			d.reconnectCount++
		}

		time.Sleep(time.Second * 5)
	}
}
