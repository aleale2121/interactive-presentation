package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func configureLogger(verbose bool, logPath string) *logrus.Logger {
	var log = logrus.New()
	if verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	formatter := new(logrus.JSONFormatter)
	formatter.TimestampFormat = "01-01-2001 13:00:00"

	log.SetFormatter(formatter)
	log.SetFormatter(&logrus.JSONFormatter{})

	var mw io.Writer
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("Unable to open log file at: %s, error: %v", logPath, err)
			os.Exit(-1)
		}

		mw = io.MultiWriter(os.Stdout, f)
	} else {
		mw = io.MultiWriter(os.Stdout)
	}

	log.SetOutput(mw)

	return log
}
