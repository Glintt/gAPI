package controllers

import (
	apianalytics "gAPIManagement/api/api-analytics"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"

	routing "github.com/qiangxue/fasthttp-routing"
)

func Logs(c *routing.Context) error {
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))

	res, status := apianalytics.AnalyticsMethods[config.GApiConfiguration.Logs.Type]["logs"].(func(string) (string, int))(apiEndpoint)

	http.Response(c, res, status, apianalytics.SERVICE_NAME, config.APPLICATION_JSON)

	return nil
}

func APIAnalytics(c *routing.Context) error {
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))

	res, status := apianalytics.AnalyticsMethods[config.GApiConfiguration.Logs.Type]["analytics"].(func(string) (string, int))(apiEndpoint)

	http.Response(c, res, status, apianalytics.SERVICE_NAME, config.APPLICATION_JSON)

	return nil
}
