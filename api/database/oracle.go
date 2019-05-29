package database

import (
	"database/sql"
	"github.com/Glintt/gAPI/api/utils"

	_ "gopkg.in/goracle.v2"
)

var dbConnection *sql.DB

func ConnectToOracle(connectionString string) (*sql.DB, error) {
	var err error

	if dbConnection != nil {
		dbConnection.Ping()
		return dbConnection, err
	}
	dbConnection, err = sql.Open("goracle", connectionString)

	if err != nil {
		return nil, err
	}

	if err != nil {
		utils.LogMessage("error connecting to oracle on "+connectionString+". Err: "+err.Error(), utils.ErrorLogType)
	} else {
		utils.LogMessage("connected to oracle on "+connectionString, utils.InfoLogType)
	}

	return dbConnection, err
}

func CloseOracleConnection(dbConnection *sql.DB) error {
	//defer dbConnection.Close()
	return nil
}
