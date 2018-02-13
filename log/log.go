package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Entry
)

func GetLogger() *logrus.Entry {
	if logger == nil {
		setup()
	}

	return logger
}

func setup() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	logger = logrus.WithFields(logrus.Fields{})
}
