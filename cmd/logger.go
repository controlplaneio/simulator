package cmd

import (
	"github.com/sirupsen/logrus"
	"os"
)

func newLogger(level string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.Out = os.Stdout

	if level == "" {
		level = "info"
	}

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.WithFields(logrus.Fields{
			"Level": level,
			"Error": err,
		}).Error("Error parsing loglevel")
	}
	log.SetLevel(parsedLevel)

	return log
}
