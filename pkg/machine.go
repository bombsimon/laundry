package pkg

import (
	"encoding/json"
	"fmt"
)

// Machine represents a laundry machine, holding an info line and a working state
type Machine struct {
	ID      int    `db:"id"      json:"id"`
	Info    string `db:"info"    json:"info"`
	Working bool   `db:"working" json:"working"`
}

// MachineHandler is the interface to implement to handle machines.
type MachineHandler interface {
	AddMachine(machine *Machine) (*Machine, error)
	GetMachine(id int) (*Machine, error)
	GetMachines() ([]*Machine, error)
	RemoveMachine(machine *Machine) error
	RemoveMachineByID(id int) error
	UpdateMachine(id int, um *Machine) (*Machine, error)
}

// UnmarshalJSON overrides the default unmarshaling to determine weather the
// value for working condition was omitted or actually passed as false
func (m *Machine) UnmarshalJSON(data []byte) error {
	var err error

	b := struct {
		Bool *bool `json:"working"`
	}{}

	if err = json.Unmarshal(data, &b); err != nil {
		return err
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
		return err
	}

	return nil
}
