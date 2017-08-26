package laundry

import (
	_ "github.com/go-sql-driver/mysql"
)

type Machine struct {
	Id      int    `db:"id" 	 json:"id"`
	Name    string `db:"name"    json:"name"`
	Working bool   `db:"working" json:"working"`
}
