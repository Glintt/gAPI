package controllers

import (
	"gAPIManagement/api/config"
	"github.com/qiangxue/fasthttp-routing"
	apianalytics "gAPIManagement/api/api-analytics"
	"gAPIManagement/api/http"
)

func Logs(c *routing.Context) error {
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))
	
	res, status := apianalytics.Logs(apiEndpoint)

	http.Response(c, res, status, apianalytics.SERVICE_NAME, config.APPLICATION_JSON)

	return nil
}

func APIAnalytics(c *routing.Context) error {
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))
	
	res, status := apianalytics.APIAnalytics(apiEndpoint)

	http.Response(c, res, status, apianalytics.SERVICE_NAME, config.APPLICATION_JSON)

	return nil
}