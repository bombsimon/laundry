package laundry

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Database       Database       `yaml:"database"`
	HTTP           Http           `yaml:"http"`
	Bookings       BookingRules   `yaml:"bookings"`
	Administration Administration `yaml:"administration"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Http struct {
	Listen string `yaml:"listen"`
}

type BookingRules struct {
	MaxAllowed      int `yaml:"max_allowed"`
	MinSlotDuration int `yaml:"min_slot_duration"`
}

type Administration struct {
	SupportEmail string `yaml:"support_email"`
}

func NewConfig(configFile string) (*Configuration, error) {
	var c Configuration

	if configFile == "" {
		configFile = "config/back-end.yaml"
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, ExtError(err).WithStatus(http.StatusInternalServerError)
	}

	if err = yaml.Unmarshal(file, &c); err != nil {
		return nil, ExtError(err).WithStatus(http.StatusInternalServerError)
	}

	return &c, nil
}
