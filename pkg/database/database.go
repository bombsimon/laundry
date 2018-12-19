package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/mysql"
)

// DBConfig represents the database configuration for the laundry service
type DBConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Database      string `yaml:"database"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	RetryCount    int    `yaml:"retry_count"`
	RetryInterval int    `yaml:"retry_interval"`
}

// SetupConnection will setup connections and store them in a
// connection pool
func SetupConnection(config DBConfig) *sql.DB {
	dsn := os.Getenv("LAUNDRY_DSN")
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=1", config.Username, config.Password, config.Host, config.Port, config.Database)
	}

	for i := config.RetryCount; i >= 0; i-- {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
			time.Sleep(time.Second * time.Duration(config.RetryInterval))
			continue
		}

		// Monitor the connection returned.
		go monitorConnection(db)

		return db
	}

	return nil
}

func monitorConnection(db *sql.DB) {
	for {
		if err := db.Ping(); err != nil {
			// Should reconnect
		}

		time.Sleep(time.Second * 5)
	}
}
