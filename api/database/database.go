package database

import (
	"gAPIManagement/api/config"
	"os"
)

var MONGO_HOST string
var MONGO_DB string
var MONGO_PORT string

var ORACLE_CONNECTION_STRING string

var IsConnectionDone = false

func InitDatabaseConnection() error {
	MONGO_HOST = os.Getenv("MONGO_HOST")
	MONGO_DB = os.Getenv("MONGO_DB")
	MONGO_PORT = os.Getenv("MONGO_PORT")
	ORACLE_CONNECTION_STRING = os.Getenv("ORACLE_CONNECTION_STRING")

	if MONGO_PORT != "" {
		MONGO_HOST = MONGO_HOST + ":" + MONGO_PORT
	}

	IsConnectionDone = true

	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" {
		err := ConnectToMongo(MONGO_HOST)

		if err != nil {
			IsConnectionDone = false
		}
		return err
	}
	return nil
}
