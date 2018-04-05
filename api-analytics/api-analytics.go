package apianalytics

import (
	"api-management/config"
	"api-management/http"
	"api-management/authentication"

	routing "github.com/qiangxue/fasthttp-routing"
)

func StartAPIAnalytics(router *routing.Router) {
	var analyticsAPI *routing.RouteGroup
	analyticsAPI = router.Group(config.ANALYTICS_GROUP)

	analyticsAPI.Get("/api",authentication.AuthorizationMiddleware, APIAnalytics)
}

func APIAnalytics(c *routing.Context) error {
	logsURL := config.ELASTICSEARCH_URL + "/request-logs-*/logs/_search"

	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))

	//skip := c.QueryArgs().Peek("skip")
	//take := c.QueryArgs().Peek("take")

	if apiEndpoint != "" {
		apiEndpoint = `"query":{
				"match":{
					"ServiceName":"` + apiEndpoint + `"
				}	
			},`
	}
	response := http.MakeRequest(config.POST, logsURL, `{
		"size": 10,
		"from": 0,
		"sort" : [
        	{ "DateTime.keyword" : {"order" : "desc"}}
    	],
		`+apiEndpoint+`
		"aggs": {
			"api": {
				"terms":{
					"field":"ServiceName.keyword",
					"order":{"_count":"desc"}
				},
				"aggs": {
					"UserAgent": {
						"terms": {
							"field": "UserAgent.keyword",
							"size": 5,
							"order": {
								"_count": "desc"
							}
						}
					},
					"StatusCode": {
						"terms": {
							"field": "StatusCode",
							"size": 5,
							"order": {
								"_count": "desc"
							}
						},
						"aggs": {
							"Count": {
								"cardinality": {
									"field": "StatusCode"
								}
							}
						}
					},
					"RemoteAddr": {
						"terms": {
							"field": "RemoteIp.keyword",
							"size": 5,
							"order": {
								"_count": "desc"
							}
						},
						"aggs": {
							"Count": {
								"cardinality": {
									"field": "RemoteIp.keyword"
								}
							}
						}
					},
					"MaxElapsedTime": {
						"max": {
							"field": "ElapsedTime"
						}
					},
					"MinElapsedTime": {
						"min": {
							"field": "ElapsedTime"
						}
					},
					"AvgElapsedTime": {
						"avg": {
							"field": "ElapsedTime"
						}
					}
				}
			}
		}
	}`, nil)

	http.Response(c, string(response.Body()), response.StatusCode(), config.ANALYTICS_GROUP+"/api")

	return nil
}
