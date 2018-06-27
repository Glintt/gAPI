package apianalytics

import (
	"strings"
	"regexp"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
)

var SERVICE_NAME = "api_analytics"

func LogsURL() string {
	return config.ELASTICSEARCH_URL + "/" + config.ELASTICSEARCH_LOGS_INDEX + "/_search"
}

func Logs(apiEndpoint string) (string, int) {	
	if apiEndpoint != "" {
		apiEndpoint = `
		"must":{
			"match_phrase" : {
				"ServiceName": "`+ apiEndpoint +`"
			}
		},`
	}

	response := http.MakeRequest(config.POST, LogsURL(), LogsQuery(apiEndpoint), nil)
		
	var re = regexp.MustCompile(`\\"Authorization\\"( )*:( )*\\"[^"]+\\"`)
	s := re.ReplaceAllString(string(response.Body()), "")
	s = strings.Replace(s, ",,", ",", -1)

	return s, response.StatusCode()
}

func APIAnalytics(apiEndpoint string) (string, int) {
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