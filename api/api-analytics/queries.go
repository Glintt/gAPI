package apianalytics

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
				`+ apiEndpoint +`
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
	}`
}