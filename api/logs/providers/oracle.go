package providers

import (
	"fmt"
	"gAPIManagement/api/config"
	"gAPIManagement/api/database"
	"gAPIManagement/api/logs/models"
	"gAPIManagement/api/utils"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const INSERT_LOG_ORACLE = `INSERT INTO gapi_request_logs(id,method,uri,request_body,host,user_agent,remote_addr,remote_ip,headers,query_args,date_time,response,elapsed_time,status_code,service_name,index_name,request_grouper_date) 
VALUES(:id,:method,:uri,:request_body,:host,:user_agent,:remote_addr,:remote_ip,:headers,:query_args,:date_time,:response,:elapsed_time,:status_code,:service_name,:index_name,:request_grouper_date)`

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
	)

	if err != nil {
		fmt.Println(dateTime)
		fmt.Println(currentDate)
		fmt.Println(err.Error())
	}

	tx.Commit()
	database.CloseOracleConnection(db)
}
