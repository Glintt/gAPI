package providers

import (
	"encoding/json"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/logs/models"
	"github.com/Glintt/gAPI/api/utils"

	"github.com/valyala/fasthttp"
)

func PublishElastic(reqLogging *models.RequestLogging) {
	utils.LogMessage("ELASTIC PUBLISH", utils.DebugLogType)
	currentDate := utils.CurrentDate()

	indexName := config.ELASTICSEARCH_LOGS_INDEX
	if reqLogging.IndexName != "" {
		indexName = reqLogging.IndexName
	}

	logsURL := config.ELASTICSEARCH_URL + "/" + indexName + "/request-logs-" + currentDate

	reqLoggingJson, _ := json.Marshal(reqLogging)

	request := fasthttp.AcquireRequest()

	request.SetRequestURI(logsURL)
	request.Header.SetMethod("POST")

	request.Header.SetContentType("application/json")
	request.SetBody(reqLoggingJson)
	client := fasthttp.Client{}

	resp := fasthttp.AcquireResponse()
	err := client.Do(request, resp)

	if err != nil {
		utils.LogMessage("ELASTIC PUBLISH - error:"+err.Error(), utils.DebugLogType)

		resp.SetStatusCode(400)
	}
}
