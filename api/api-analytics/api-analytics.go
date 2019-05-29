package apianalytics

import (
	"github.com/Glintt/gAPI/api/api-analytics/providers"
)

const SERVICE_NAME = "api_analytics"

const (
	LogElasticType     = "Elastic"
	LogOracleType      = "Oracle"
	LogRabbitQueueType = "Rabbit"
)

var AnalyticsMethods = map[string]map[string]interface{}{
	"Elastic": {
		"logs":          providers.LogsElastic,
		"analytics":     providers.APIAnalyticsElastic,
		"app_analytics": providers.ApplicationAnalyticsElastic,
	},

	"Oracle": {
		"logs":          providers.LogsOracle,
		"analytics":     providers.APIAnalyticsOracle,
		"app_analytics": providers.ApplicationAnalyticsOracle,
	},
}
