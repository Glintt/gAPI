package utils

import (
	
)

func PreventCrash(){
	if r := recover(); r != nil {
		LogMessage("Publish Log Panic recover", DebugLogType)
	}
}