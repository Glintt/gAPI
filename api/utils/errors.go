package utils

import (
	"fmt"
)

func PreventCrash(){
	if r := recover(); r != nil {
		fmt.Println("Publish Log Panic", r)
	}
}