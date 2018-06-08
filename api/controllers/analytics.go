package controllers

import (
	"gAPIManagement/api/config"
	"github.com/qiangxue/fasthttp-routing"
	"gAPIManagement/api/api-analytics"
	"gAPIManagement/api/http"
)

func Logs(c *routing.Context) error {
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))
	
	res, status := apianalytics.Logs(apiEndpoint)

	http.Response(c, res, status, config.ANALYTICS_GROUP+"/logs")

	return nil
}

func APIAnalytics(c *routing.Context) error {
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))
	
	res, status := apianalytics.APIAnalytics(apiEndpoint)

	http.Response(c, res, status, config.ANALYTICS_GROUP+"/api")

	return nil
}