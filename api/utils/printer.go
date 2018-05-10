package utils

import (
	"time"
	"fmt"
)

func LogMessage(message string) {
	fmt.Print(time.UnixDate + " - ")

	fmt.Println(message)
}