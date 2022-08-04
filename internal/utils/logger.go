package utils

import (
	"github.com/sirupsen/logrus"
	"io"
)

func DefaultLogger(out io.Writer, level logrus.Level) *logrus.Entry {
	log := logrus.New()
	log.SetOutput(out)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Level = level

	return logrus.NewEntry(log)
}
