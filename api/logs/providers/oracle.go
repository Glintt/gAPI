package providers

import (
	"time"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/logs/models"
	"github.com/Glintt/gAPI/api/utils"

	"github.com/robfig/cron"
	"gopkg.in/mgo.v2/bson"
)

const INSERT_LOG_ORACLE = `INSERT INTO gapi_request_logs(id,method,uri,request_body,host,user_agent,remote_addr,remote_ip,headers,query_args,date_time,response,elapsed_time,status_code,service_name,index_name,request_grouper_date,other_info) 
VALUES(:id,:method,:uri,:request_body,:host,:user_agent,:remote_addr,:remote_ip,:headers,:query_args,:date_time,:response,:elapsed_time,:status_code,:service_name,:index_name,:request_grouper_date,:other_info)`
const DELETE_OLDER_LOGS_SUCCESS = `DELETE FROM gapi_request_logs
WHERE     date_time <=
				SYSDATE
			  - (SELECT config_value
				   FROM gapi_configurations
				  WHERE config_key =
							'GAPI_REQUEST_LOGS_SUCCESS_DURATION_DAYS')
	  AND status_code < 300`
const DELETE_OLDER_LOGS_ERROR = `DELETE FROM gapi_request_logs
	  WHERE     date_time <=
					  SYSDATE
					- (SELECT config_value
						 FROM gapi_configurations
						WHERE config_key =
								  'GAPI_REQUEST_LOGS_ERRORS_DURATION_DAYS')
			AND status_code >= 300`

const dateFormat = "2-Jan-06 3:04:05.000000 PM"

func RemoveOldLogsOracle() {
	c := cron.New()
	c.AddFunc("@daily", func() {
		db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
		if err != nil {
			utils.LogMessage("Remove old logs database connection failed: "+err.Error(), utils.InfoLogType)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			utils.LogMessage("Remove old logs database connection Transaction creation failed: "+err.Error(), utils.InfoLogType)
			return
		}
		tx.Exec(DELETE_OLDER_LOGS_ERROR)
		tx.Exec(DELETE_OLDER_LOGS_SUCCESS)
		tx.Commit()
		database.CloseOracleConnection(db)
		utils.LogMessage("Removed old logs oracle job run ", utils.InfoLogType)

	})
	c.Start()
	utils.LogMessage("Remove old logs oracle job started ", utils.InfoLogType)

}

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
