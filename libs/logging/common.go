package logging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var lastLogsFileDate *time.Time
var logsFile *os.File
var log = logrus.New()
var logsDirectoryPath = "./logs/"

func Init() {
	createDirIfNotExist(logsDirectoryPath)
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})
}

func Log() *logrus.Logger {
	setOutputFile()

	return log
}

func setOutputFile() {
	filename := logsDirectoryPath + getCurrentLogsFilename()
	if logsFile == nil || logsFile.Name() != filename {
		if logsFile != nil {
			err := logsFile.Close()
			if err != nil {
				log.Error("can't close logs file: %s", err)
			}
		}

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			// Output to stdout instead of the default stderr
			log.SetOutput(file)
			logsFile = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}
}

func getCurrentLogsFilename() string {
	currentTime := time.Now()
	if lastLogsFileDate == nil {
		lastLogsFileDate = &currentTime
	} else {
		if currentTime.Day() > lastLogsFileDate.Day() || currentTime.Month() > lastLogsFileDate.Month() || currentTime.Year() > lastLogsFileDate.Year() {
			lastLogsFileDate = &currentTime
		}
	}

	return "logs_" + lastLogsFileDate.Format("2006_01_02") + ".log"
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
