package database

import (
	"os"
)

var MONGO_HOST string
var MONGO_DB string
var MONGO_PORT string

var IsConnectionDone = false

func InitDatabaseConnection() error {
	MONGO_HOST = os.Getenv("MONGO_HOST")
	MONGO_DB = os.Getenv("MONGO_DB")
	MONGO_PORT = os.Getenv("MONGO_PORT")

	if MONGO_PORT != "" {
		MONGO_HOST = MONGO_HOST + ":" + MONGO_PORT
	}
	err := ConnectToMongo(MONGO_HOST)
	
	if err == nil {
		IsConnectionDone = true
	}

	return err
}