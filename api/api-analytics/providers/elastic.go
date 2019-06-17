package providers

import (
	"regexp"
	"strings"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery/appgroups"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
)

func LogsURL() string {
	return config.ELASTICSEARCH_URL + "/" + config.ELASTICSEARCH_LOGS_INDEX + "/_search"
}

func LogsElastic(apiEndpoint string) (string, int) {
	if apiEndpoint != "" {
		apiEndpoint = `
		"must":{
			"match_phrase" : {
				"ServiceName": "` + apiEndpoint + `"
			}
		},`
	}

	response := http.MakeRequest(config.POST, LogsURL(), LogsQuery(apiEndpoint), nil)

	var re = regexp.MustCompile(`\\"Authorization\\"( )*:( )*\\"[^"]+\\"`)
	s := re.ReplaceAllString(string(response.Body()), "")
	s = strings.Replace(s, ",,", ",", -1)

	return s, response.StatusCode()
}

func APIAnalyticsElastic(apiEndpoint string) (string, int) {
	if apiEndpoint != "" {
		apiEndpoint = `"query":{
				"match_phrase":{
					"ServiceName":"` + apiEndpoint + `"
				}	
			},`
	}

	response := http.MakeRequest(config.POST, LogsURL(), APIAnalyticsQuery(apiEndpoint), nil)

	return string(response.Body()), response.StatusCode()
}

func ApplicationAnalyticsElastic(appGroupId string) (string, int) {
	appGroup, _ := appgroups.ApplicationGroupMethods[database.SD_TYPE]["getbyid"].(func(string) (appgroups.ApplicationGroup, error))(appGroupId)
	appGroupServices, _ := appgroups.ApplicationGroupMethods[database.SD_TYPE]["getservicesforappgroup"].(func(appgroups.ApplicationGroup) ([]service.Service, error))(appGroup)

	response := http.MakeRequest(config.POST, LogsURL(), ApplicationAnalyticsQuery(appGroupServices), nil)

	return string(response.Body()), response.StatusCode()
}

func LogsQuery(apiEndpoint string) string {
	return `
	{
		"size": 10,
		"from": 0,
		"sort" : [
			{ "DateTime.keyword" : {"order" : "desc"}}
		],
		"query":{
			"bool": {
				` + apiEndpoint + `
				"filter" : {
					"range" : {
						"StatusCode" : {
							"gte" : 300
						}
					}
				}
			}
			
		}
	}`
}

func APIAnalyticsQuery(apiEndpoint string) string {
	return `{
		"size": 0,
		"from": 0,
		` + apiEndpoint + `
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
	}`
}

func ApplicationAnalyticsQuery(services []service.Service) string {
	shouldQuery := ""
	for _, s := range services {
		if shouldQuery != "" {
			shouldQuery = shouldQuery + ","
		}
		shouldQuery = shouldQuery + `{ "match": { "ServiceName": "` + s.MatchingURI + `" }}`
	}
	return `{
		"size": 0,
		"from": 0,
		"query":{
			"bool":{
				"should": [` + shouldQuery + `
				]
				
			}
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
	`
}
