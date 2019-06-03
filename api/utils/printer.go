package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron"
)

const (
	FILE_LOGS_ENV_VAR     = "LOGS_TO_FILE"
	DebugLogType          = "DEBUG"
	InfoLogType           = "INFO"
	ErrorLogType          = "ERROR"
	logDateFilenameFormat = "2006-01-02"
	serverLogsFolder      = "gapi_log_files"
	logFileExtension      = "log"
)

var initialized = false

// Logs duration (15 days)
var logsPersistenceTime = 7 * 24 * time.Hour

// Remove files older than
func deletefiles(path string, f os.FileInfo, err error) (e error) {
	logsPath := "." + string(filepath.Separator) + serverLogsFolder
	now := time.Now()
	fileInfo, err := ioutil.ReadDir("." + string(filepath.Separator) + serverLogsFolder)
	for _, info := range fileInfo {
		if diff := now.Sub(info.ModTime()); diff > logsPersistenceTime {
			pathToDelete := logsPath + string(filepath.Separator) + info.Name()
			LogMessage("Deleting "+pathToDelete, DebugLogType)
			err = os.Remove(pathToDelete)
			if err != nil {
				LogMessage("Error deleting "+pathToDelete+" with error: "+err.Error(), DebugLogType)
			}
		}
	}
	return
}

// Configure everything required when using files to log
func configureLoggingFiles() {
	if initialized {
		return
	}

	path := "." + string(filepath.Separator) + serverLogsFolder
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	c := cron.New()
	c.AddFunc("* * * * *", func() {
		filepath.Walk("."+string(filepath.Separator)+serverLogsFolder, deletefiles)
	})
	c.Start()
	initialized = true
}

func LogMessage(message string, logtype string) {
	if logtype == DebugLogType && os.Getenv("DEBUG") != "true" {
		return
	}

	currentTime := time.Now()
	currDate := currentTime.UTC().String()
	var logger = log.New(os.Stdout, currDate+" - ", log.LstdFlags)

	// If log to file enabled, then change output to file inside gapi_log_files folder
	if os.Getenv(FILE_LOGS_ENV_VAR) != "" && os.Getenv(FILE_LOGS_ENV_VAR) == "true" {
		configureLoggingFiles()
		logFileName := "." + string(filepath.Separator) + serverLogsFolder + string(filepath.Separator) + currentTime.Format(logDateFilenameFormat) + "." + logFileExtension
		f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			logger.SetOutput(f)
		} else {
			logger.Println(err.Error())
		}
		defer f.Close()
	}

	logger.Println(message)
}
