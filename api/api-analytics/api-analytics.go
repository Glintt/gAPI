package apianalytics

import (
	"strings"
	"regexp"
	"fmt"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/authentication"

	routing "github.com/qiangxue/fasthttp-routing"
)
var logsURL string


func StartAPIAnalytics(router *routing.Router) {
	logsURL = config.ELASTICSEARCH_URL + "/request-logs-*/logs/_search"
	var analyticsAPI *routing.RouteGroup
	analyticsAPI = router.Group(config.ANALYTICS_GROUP)

	analyticsAPI.Get("/api", authentication.AuthorizationMiddleware, APIAnalytics)
	analyticsAPI.Get("/logs", authentication.AuthorizationMiddleware, Logs)
}

func Logs(c *routing.Context) error {
	fmt.Println("LOGS="+logsURL)
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))
	if apiEndpoint != "" {
		apiEndpoint = `
		"must":{
			"match" : {
				"ServiceName": "`+ apiEndpoint +`"
			}
		},`
	}
	
	response := http.MakeRequest(config.POST, logsURL, `
		{
			"size": 10,
			"from": 0,
			"sort" : [
				{ "DateTime.keyword" : {"order" : "desc"}}
			],
			"query":{
				"bool": {
					`+ apiEndpoint +`
					"must_not" : {
						"range" : {
							"StatusCode" : {
								"lt" : 300
							}
						}
					}
				}
				
			}
		}`, nil)
		
		
	var re = regexp.MustCompile(`\\"Authorization\\":\\"[\w\d\s.-]+\\"`)
	s := re.ReplaceAllString(string(response.Body()), "")
	s = strings.Replace(s, ",,", ",", -1)
		
	http.Response(c, string(response.Body()), response.StatusCode(), config.ANALYTICS_GROUP+"/logs")

	return nil
}

func APIAnalytics(c *routing.Context) error {
	fmt.Println(logsURL)
	apiEndpoint := string(c.QueryArgs().Peek("endpoint"))

	if apiEndpoint != "" {
		apiEndpoint = `"query":{
				"match":{
					"ServiceName":"` + apiEndpoint + `"
				}	
			},`
	}
	response := http.MakeRequest(config.POST, logsURL, `{
		"size": 0,
		"from": 0,
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
