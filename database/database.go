package database

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bombsimon/laundry/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	connectionPool []dbConnection
)

type dbConnection struct {
	db             *sqlx.DB
	reconnected    bool
	reconnectCount int
}

// SetupConnection will setup connections and store them in a
// connection pool
func SetupConnection(config config.Database) {
	dsn := os.Getenv("LAUNDRY_DSN")
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=1", config.Username, config.Password, config.Host, config.Port, config.Database)
	}

	for i := 1; i <= config.PoolSize; i++ {
		var dbConn dbConnection

		for j := config.RetryCount; j >= 0; j-- {
			db, err := sqlx.Connect("mysql", dsn)

			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}

			dbConn.db = db
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

func monitorConnection(d dbConnection) {
	for {
		if err := d.db.Ping(); err != nil {
			// Should reconnect
			d.reconnected = true
			d.reconnectCount += 1
		}

		time.Sleep(time.Second * 5)
	}
}
