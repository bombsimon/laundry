package laundry

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bombsimon/laundry/database"
	"github.com/bombsimon/laundry/errors"
	_ "github.com/go-sql-driver/mysql"
)

// Machine represents a laundry machine, holding an info line and a working state
type Machine struct {
	Id      int    `db:"id"      json:"id"`
	Info    string `db:"info"    json:"info"`
	Working bool   `db:"working" json:"working"`
}

// UnmarshalJSON overrides the default unmarshaling to determine weather the
// value for working condition was ommitted or actually passed as false
func (m *Machine) UnmarshalJSON(data []byte) *errors.LaundryError {
	var err error

	b := struct {
		Bool *bool `json:"working"`
	}{}

	if err = json.Unmarshal(data, &b); err != nil {
		return errors.New(err)
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

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// GetMachines returns a list of all Machines added in the database. If there are
// no machines, an empty list will be returned. The same applies if an error occurs.
func GetMachines() ([]Machine, *errors.LaundryError) {
	db := database.GetConnection()

	var machines []Machine

	rows, err := db.Queryx("SELECT * FROM machines")
	if err != nil {
		return machines, errors.New("Could not get machines").CausedBy(err)
	}

	defer rows.Close()

	for rows.Next() {
		var m Machine
		if err := rows.StructScan(&m); err != nil {
			return machines, errors.New("Could not get machines").CausedBy(err)
		}

		machines = append(machines, m)
	}

	return machines, nil
}

// GetMachine will return the Machine with passed ID if it exists in the database.
// If an error occurs or the Machine does not exist, nil will be returned.
func GetMachine(id int) (*Machine, *errors.LaundryError) {
	db := database.GetConnection()

	var m Machine
	if err := db.QueryRowx("SELECT * FROM machines WHERE id = ?", id).StructScan(&m); err == sql.ErrNoRows {
		return nil, errors.New("Machine with id %d not found", id).WithStatus(http.StatusNotFound)
	} else if err != nil {
		return nil, errors.New("Could not get row").CausedBy(err)
	}

	return &m, nil
}

// AddMachine will take a defined Machine and add it in the database. If the Machine
// has an id or is an existing Machine, the id will be omitted and a copy will be created.
func AddMachine(m *Machine) (*Machine, *errors.LaundryError) {
	db := database.GetConnection()

	if m.Info == "" {
		return nil, errors.New("Missing info in request")
	}

	row, err := db.Exec("INSERET INTO machines ( info, working ) VALUES ( ?, ? )", m.Info, m.Working)
	if err != nil {
		return nil, errors.New("Could not create machine").CausedBy(err)
	}

	lastId, err := row.LastInsertId()
	if err != nil {
		return nil, errors.New(err)
	}

	m.Id = int(lastId)

	return m, nil
}

// UpdateMachine will take a machine and update the row with the corresponding id.
// The passed Machine will be returned if successful.
func UpdateMachine(id int, m *Machine) (*Machine, *errors.LaundryError) {
	db := database.GetConnection()

	machine, lErr := GetMachine(id)
	if lErr != nil {
		return nil, lErr
	}

	if m.Info == "" {
		return nil, errors.New("Missing field info").WithStatus(http.StatusBadRequest)
	}

	machine.Info = m.Info
	machine.Working = m.Working

	if _, err := db.Exec("UPDATE machines SET info = ?, working = ? WHERE id = ?", machine.Info, machine.Working, machine.Id); err != nil {
		return nil, errors.New("Could not update machine with id %d", m.Id).CausedBy(err)
	}

	return m, nil
}

// RemoveMachine will remove a machine alltogether. If the Machine is related
// to any slots in the booking system that will be removed aswell
func RemoveMachine(m *Machine) *errors.LaundryError {
	db := database.GetConnection()

	if _, err := db.Exec("DELETE FROM machines WHERE ID = ?", m.Id); err != nil {
		return errors.New("Could not remove machine with id %d", m.Id).CausedBy(err)
	}

	return nil
}

// RemoveMachineById will remove a machine by sending the machine id.
func RemoveMachineById(id int) *errors.LaundryError {
	m, err := GetMachine(id)
	if err != nil {
		return err
	}

	return RemoveMachine(m)
}
