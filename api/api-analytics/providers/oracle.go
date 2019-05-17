package providers

import (
	"database/sql"
	"encoding/json"
	"gAPIManagement/api/database"
	logsModels "gAPIManagement/api/logs/models"
	"strings"
)

const LOGS_QUERY_ORACLE = `SELECT id,method,uri,request_body,host,user_agent,remote_addr,remote_ip,headers,query_args,date_time,response,elapsed_time,status_code,service_name,index_name,request_grouper_date FROM gapi_request_logs
where status_code >= 300
`

func LogsOracle(apiEndpoint string) (string, int) {
	var err error
	var rows *sql.Rows

	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + ` "}`, 500
	}

	query := LOGS_QUERY_ORACLE

	if apiEndpoint != "" {
		query = strings.Replace(query, "#WHERE_CLAUSE#", " where service_name = :serviceName", -1)
	}
	query = `SELECT * FROM
	(
		SELECT a.*, rownum r__
		FROM
		(
			` + query + `
			) a
			WHERE rownum < ((:page * 10) + 1 )
		)
		WHERE r__ >= (((:page-1) * 10) + 1)`

	if apiEndpoint != "" {
		rows, err = db.Query(query, apiEndpoint, 1)
	} else {
		rows, err = db.Query(query, 1)
	}
	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + ` "}`, 500
	}
	requestLogs := RowsToLogRequestModel(rows, true)

	defer rows.Close()
	database.CloseOracleConnection(db)

	var res []map[string]interface{}
	for _, v := range requestLogs {
		e := make(map[string]interface{})
		e["_source"] = v
		e["_id"] = v.Id
		res = append(res, e)
	}
	respString, _ := json.Marshal(res)
	return `{"hits":{"hits": ` + string(respString) + ` }}`, 200
}

func APIAnalyticsOracle(apiEndpoint string) (string, int) {
	return "", 500
}

func RowsToLogRequestModel(rows *sql.Rows, containsPagination bool) []logsModels.RequestLogging {
	var logs []logsModels.RequestLogging
	for rows.Next() {
		var s logsModels.RequestLogging
		var a int
		var currentDate string
		if containsPagination {
			rows.Scan(
				&s.Id,
				&s.Method,
				&s.Uri,
				&s.RequestBody,
				&s.Host,
				&s.UserAgent,
				&s.RemoteAddr,
				&s.RemoteIp,
				&s.Headers,
				&s.QueryArgs,
				&s.DateTime,
				&s.Response,
				&s.ElapsedTime,
				&s.StatusCode,
				&s.ServiceName,
				&s.IndexName,
				&currentDate,
				&a,
			)
		} else {
			rows.Scan(
				&s.Id, &s.Method,
				&s.Uri,
				&s.RequestBody,
				&s.Host,
				&s.UserAgent,
				&s.RemoteAddr,
				&s.RemoteIp,
				&s.Headers,
				&s.QueryArgs,
				&s.DateTime,
				&s.Response,
				&s.ElapsedTime,
				&s.StatusCode,
				&s.ServiceName,
				&s.IndexName,
				&currentDate,
			)
		}

		logs = append(logs, s)
	}

	defer rows.Close()

	return logs
}
