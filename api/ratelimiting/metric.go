package ratelimiting

import (
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/servicediscovery/service"

	routing "github.com/qiangxue/fasthttp-routing"
)

func GetRateLimitingMetricName(c *routing.Context, limiter config.GApiRateLimitingConfig) string {
	metricName := ""
	for _, val := range limiter.Metrics {
		if metricName != "" {
			metricName = metricName + "-"
		}
		switch val {
		case "RemoteAddr":
			metricName = metricName + remoteAddr(c)
			break
		case "MatchingUri":
			metricName = metricName + matchingUri(c)
			break

		default:
			metricName = metricName + val
			break
		}
	}
	return metricName
}

func remoteAddr(c *routing.Context) string {
	return c.RemoteAddr().String()
}

func matchingUri(c *routing.Context) string {
	return serviceForUri(c).MatchingURI
}

func serviceForUri(c *routing.Context) service.Service {
	s, _ := sd.GetEndpointForUri(string(c.Request.RequestURI()))
	return s
}
