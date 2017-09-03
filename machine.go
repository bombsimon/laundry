package laundry

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Machine represents a laundry machine, holding an info line and a working state
type Machine struct {
	Id      int    `db:"id"      json:"-"`
	Info    string `db:"info"    json:"info"`
	Working bool   `db:"working" json:"working"`
}

// UnmarshalJSON overrides the default unmarshaling to determine weather the
// value for working condition was ommitted or actually passed as false
func (m *Machine) UnmarshalJSON(data []byte) (err error) {
	b := struct {
		Bool *bool `json:"working"`
	}{}

	if err = json.Unmarshal(data, &b); err != nil {
		return
	}

	switch b.Bool {
	case nil:
		err = fmt.Errorf("Missing parameter working")
	default:
		m1 := struct {
			Info    string
			Working bool
		}{}
		err = json.Unmarshal(data, &m1)

		m.Info = m1.Info
		m.Working = m1.Working
	}

	return
}

// GetMachines returns a list of all Machines added in the database. If there are
// no machines, an empty list will be returned. The same applies if an error occurs.
func (l *Laundry) GetMachines() ([]Machine, error) {
	sqlStmt := `SELECT * FROM machines`

	var machines []Machine

	rows, err := l.db.Queryx(sqlStmt)
	if err != nil {
		l.Logger.Errorf("Could not get bookers: %s", err)
		return machines, ExtError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var m Machine
		if err := rows.StructScan(&m); err != nil {
			l.Logger.Errorf("Could not get machines: %s", err)
			return machines, ExtError(err)
		}

		machines = append(machines, m)
	}

	return machines, nil
}

// GetMachine will return the Machine with passed ID if it exists in the database.
// If an error occurs or the Machine does not exist, nil will be returned.
func (l *Laundry) GetMachine(id int) (*Machine, error) {
	sqlStmt := `SELECT * FROM machines WHERE id = ?`
	stmt, err := l.db.Preparex(sqlStmt)

	defer stmt.Close()

	if err != nil {
		l.Logger.Errorf("Could not prepare statement: %s", err)
		return nil, ExtError(err)
	}

	var m Machine
	if err = stmt.QueryRowx(id).StructScan(&m); err == sql.ErrNoRows {
		l.Logger.Warnf("Machine with ID %d not found", id)
		return nil, NewError(fmt.Sprintf("Machine with id %d not found", id)).WithStatus(404)
	} else if err != nil {
		l.Logger.Errorf("Could not get row: %s", err)
		return nil, ExtError(err)
	}

	return &m, nil
}

// AddMachine will take a defined Machine and add it in the database. If the Machine
// has an id or is an existing Machine, the id will be omitted and a copy will be created.
func (l *Laundry) AddMachine(m *Machine) (*Machine, error) {
	sqlStmt := `INSERT INTO machines ( info, working ) VALUES ( ?, ? )`

	stmt, err := l.db.Preparex(sqlStmt)

	defer stmt.Close()

	if err != nil {
		l.Logger.Errorf("Could not prepare statement: %s", err)
		return nil, ExtError(err)
	}

	row, err := stmt.Exec(m.Info, m.Working)
	if err != nil {
		l.Logger.Errorf("Could not create machine: %s", err)
		return nil, ExtError(err)
	}

	lastId, err := row.LastInsertId()
	if err != nil {
		return nil, ExtError(err)
	}

	m.Id = int(lastId)

	return m, nil
}

// UpdateMachine will take a machine and update the row with the corresponding id.
// The passed Machine will be returned if successful.
func (l *Laundry) UpdateMachine(m *Machine) (*Machine, error) {
	sqlStmt := `UPDATE machines SET info = ?, working = ? WHERE id = ?`

	stmt, err := l.db.Preparex(sqlStmt)
	if err != nil {
		l.Logger.Errorf("Could not prepare statement: %s", err)
		return nil, ExtError(err)
	}

	if _, err = stmt.Exec(m.Info, m.Working, m.Id); err != nil {
		l.Logger.Errorf("Could not update machine with id %d: %s", m.Id, err)
		return nil, ExtError(err)
	}

	return m, nil
}

// RemoveMachine will remove a machine alltogether. If the Machine is related
// to any slots in the booking system that will be removed aswell
func (l *Laundry) RemoveMachine(m *Machine) error {
	sqlStmt := `DELETE FROM machines WHERE id = ?`

	stmt, err := l.db.Preparex(sqlStmt)
	if err != nil {
		l.Logger.Errorf("Could not prepare statement: %s", err)
		return ExtError(err)
	}

	if _, err := stmt.Exec(m.Id); err != nil {
		l.Logger.Errorf("Could not remove machine with id %d: %s", m.Id, err)
		return ExtError(err)
	}

	return nil
}
