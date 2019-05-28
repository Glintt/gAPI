package controllers

import (
	apianalytics "github.com/Glintt/gAPI/api/api-analytics"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"

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

func ApplicationGroupAnalytics(c *routing.Context) error {
	applicationGroup := string(c.QueryArgs().Peek("app_id"))

	res, status := apianalytics.AnalyticsMethods[config.GApiConfiguration.Logs.Type]["app_analytics"].(func(string) (string, int))(applicationGroup)

	http.Response(c, res, status, apianalytics.SERVICE_NAME, config.APPLICATION_JSON)

	return nil
}
