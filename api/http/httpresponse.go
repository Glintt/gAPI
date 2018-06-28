package http

import (
	"gAPIManagement/api/utils"
	
	routing "github.com/qiangxue/fasthttp-routing"
)

type ResponseInfo struct {
	StatusCode  int
	ContentType []byte
	Body        []byte
}

func Response(c *routing.Context, response string, statuscode int, service string, contentType string) {

	utils.LogMessage("RESPONSE ==> " + response, utils.InfoLogType)

	c.Response.SetBody([]byte(response))
	c.Response.Header.SetContentType(contentType)
	c.Response.Header.Set("service", service)
	c.Response.Header.SetStatusCode(statuscode)
/* 
	if ! config.GApiConfiguration.Logs.Active {
		return
	}

	elapsedTime := 10
	serviceElapsedTime := 10

	queryArgs, _ := json.Marshal(GetQueryParamsFromRequest(c))
	headers, _ := json.Marshal(GetHeadersFromRequest(c.Request))
	logRequest := logs.NewRequestLogging(c, queryArgs, headers, utils.CurrentDateWithFormat(time.UnixDate), int64(elapsedTime), service, int64(serviceElapsedTime))
	work := logs.LogWorkRequest{Name: "", LogToSave: logRequest}
	logs.WorkQueue <- work */

}
