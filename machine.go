package laundry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bombsimon/laundry/database"
	"github.com/bombsimon/laundry/errors"
	// MySQL driver for goqu
	goqu "gopkg.in/doug-martin/goqu.v4"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
)

// Machine represents a laundry machine, holding an info line and a working state
type Machine struct {
	ID      int    `db:"id"      json:"id"`
	Info    string `db:"info"    json:"info"`
	Working bool   `db:"working" json:"working"`
}

// UnmarshalJSON overrides the default unmarshaling to determine weather the
// value for working condition was ommitted or actually passed as false
func (m *Machine) UnmarshalJSON(data []byte) error {
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
	db := database.GetGoqu()

	var machines []Machine
	if err := db.From("machines").ScanStructs(&machines); err != nil {
		return machines, errors.New("Could not get machines").CausedBy(err)
	}

	return machines, nil
}

// GetMachine will return the Machine with passed ID if it exists in the database.
// If an error occurs or the Machine does not exist, nil will be returned.
func GetMachine(id int) (*Machine, *errors.LaundryError) {
	db := database.GetGoqu()

	var m Machine
	found, err := db.From("machines").Where(goqu.Ex{
		"id": id,
	}).ScanStruct(&m)

	if err != nil {
		return nil, errors.New("Could not get row").CausedBy(err)
	}

	if !found {
		return nil, errors.New("Machine with id %d not found", id).WithStatus(http.StatusNotFound)
	}

	return &m, nil
}

// AddMachine will take a defined Machine and add it in the database. If the Machine
// has an id or is an existing Machine, the id will be omitted and a copy will be created.
func AddMachine(m *Machine) (*Machine, *errors.LaundryError) {
	if m.Info == "" {
		return nil, errors.New("Missing info in request")
	}

	db := database.GetGoqu()

	insert := db.From("machines").Insert(goqu.Record{
		"info":    m.Info,
		"working": m.Working,
	})

	row, err := insert.Exec()
	if err != nil {
		return nil, errors.New("Could not create machine").CausedBy(err)
	}

	lastID, err := row.LastInsertId()
	if err != nil {
		return nil, errors.New(err)
	}

	m.ID = int(lastID)

	return m, nil
}

// UpdateMachine will take a machine and update the row with the corresponding id.
// The passed Machine will be returned if successful.
func UpdateMachine(id int, um *Machine) (*Machine, *errors.LaundryError) {
	m, lErr := GetMachine(id)
	if lErr != nil {
		return nil, lErr
	}

	if um.Info == "" {
		return nil, errors.New("Missing field info").WithStatus(http.StatusBadRequest)
	}

	m.Info = um.Info
	m.Working = um.Working

	db := database.GetGoqu()

	update := db.From("machines").
		Where(goqu.Ex{
			"id": m.ID,
		}).
		Update(goqu.Record{
			"info":    m.Info,
			"working": m.Working,
		})

	if _, err := update.Exec(); err != nil {
		return nil, errors.New("Could not update machine with id %d", m.ID).CausedBy(err)
	}

	return m, nil
}

// RemoveMachine will remove a machine alltogether. If the Machine is related
// to any slots in the booking system that will be removed aswell
func RemoveMachine(m *Machine) *errors.LaundryError {
	db := database.GetGoqu()

	delete := db.From("machines").
		Where(goqu.Ex{
			"id": m.ID,
		}).
		Delete()

	if _, err := delete.Exec(); err != nil {
		return errors.New("Could not remove machine with id %d", m.ID).CausedBy(err)
	}

	return nil
}

// RemoveMachineByID will remove a machine by sending the machine id.
func RemoveMachineByID(id int) *errors.LaundryError {
	m, err := GetMachine(id)
	if err != nil {
		return err
	}

	return RemoveMachine(m)
}
