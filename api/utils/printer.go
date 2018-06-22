package utils

import (
	"os"
	"time"
	"fmt"
)

var ErrorLogType = "ERROR"
var InfoLogType = "INFO"
var DebugLogType = "DEBUG"

func LogMessage(message string, logtype  string) {
    if logtype == DebugLogType && os.Getenv("DEBUG") != "true" {
        return
    }
	fmt.Print(time.Now().UTC().String()  + " - ")

	fmt.Println(message)
}