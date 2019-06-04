package providers

import (
	"time"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/logs/models"
	"github.com/Glintt/gAPI/api/utils"

	"gopkg.in/mgo.v2/bson"
)

const INSERT_LOG_ORACLE = `INSERT INTO gapi_request_logs(id,method,uri,request_body,host,user_agent,remote_addr,remote_ip,headers,query_args,date_time,response,elapsed_time,status_code,service_name,index_name,request_grouper_date,other_info) 
VALUES(:id,:method,:uri,:request_body,:host,:user_agent,:remote_addr,:remote_ip,:headers,:query_args,:date_time,:response,:elapsed_time,:status_code,:service_name,:index_name,:request_grouper_date,:other_info)`

const dateFormat = "2-Jan-06 3:04:05.000000 PM"

func PublishOracle(reqLogging *models.RequestLogging) {
	utils.LogMessage("ORACLE PUBLISH", utils.DebugLogType)

	currentDate := utils.CurrentDateWithFormat(dateFormat)
	indexName := config.ELASTICSEARCH_LOGS_INDEX
	if reqLogging.IndexName != "" {
		indexName = reqLogging.IndexName
	}

	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return
	}

	dateTime, _ := time.Parse(time.UnixDate, reqLogging.DateTime)

	tx, err := db.Begin()
	if err != nil {
		utils.LogMessage("Transaction creation failed: "+err.Error(), utils.DebugLogType)
		return
	}

	_, err = tx.Exec(INSERT_LOG_ORACLE,
		bson.NewObjectId().Hex(),
		reqLogging.Method,
		reqLogging.Uri,
		reqLogging.RequestBody,
		reqLogging.Host,
		reqLogging.UserAgent,
		reqLogging.RemoteAddr,
		reqLogging.RemoteIp,
		reqLogging.Headers,
		reqLogging.QueryArgs,
		dateTime.Format(dateFormat),
		reqLogging.Response,
		reqLogging.ElapsedTime,
		reqLogging.StatusCode,
		reqLogging.ServiceName,
		indexName,
		currentDate,
		reqLogging.GetOtherInfo(),
	)

	if err != nil {
		utils.LogMessage("ORACLE LOGS PUBLISH ERROR - "+err.Error(), utils.DebugLogType)
		tx.Rollback()
	} else {
		tx.Commit()
	}

	database.CloseOracleConnection(db)
}
