package database

import (
	"database/sql"
	"gAPIManagement/api/utils"

	_ "gopkg.in/goracle.v2"
)

func ConnectToOracle(connectionString string) (*sql.DB, error) {
	var err error

	dbConnection, err := sql.Open("goracle", connectionString)

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
	defer dbConnection.Close()
	return nil
}
