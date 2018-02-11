package config

import (
	"io/ioutil"
	"net/http"

	"github.com/bombsimon/laundry/errors"
	yaml "gopkg.in/yaml.v2"
)

// Configuration represents the full configuration for the laundry service
// and the laundry RESTful API
type Configuration struct {
	Database       Database       `yaml:"database"`
	HTTP           Http           `yaml:"http"`
	Bookings       BookingRules   `yaml:"bookings"`
	Administration Administration `yaml:"administration"`
}

// Database represents the database configuration for the laundry service
type Database struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Database      string `yaml:"database"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	RetryCount    int    `yaml:"retry_count"`
	RetryInterval int    `yaml:"retry_interval"`
	PoolSize      int    `yaml:"pool_size"`
}

// Http represents the HTTP configuration for the laundry RESTful API
type Http struct {
	Listen string `yaml:"listen"`
}

// BookingRules represents the rules to be used in the laundry service
type BookingRules struct {
	MaxAllowed      int `yaml:"max_allowed"`
	MinSlotDuration int `yaml:"min_slot_duration"`
}

// Administration represents administration information for the laundry service
type Administration struct {
	SupportEmail string `yaml:"support_email"`
}

// New will create a new configuration based on a YAML file.
// The argument passed to New() is the path to a YAML file.
func New(configFile string) (*Configuration, error) {
	var c Configuration

	if configFile == "" {
		configFile = "files/back-end.yaml"
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New(err).WithStatus(http.StatusInternalServerError)
	}

	if err = yaml.Unmarshal(file, &c); err != nil {
		return nil, errors.New(err).WithStatus(http.StatusInternalServerError)
	}

	return &c, nil
}