package config

import (
	"os"
)

var SERVICE_DISCOVERY_URL = "http://localhost:8080"
var SERVICE_DISCOVERY_GROUP = "/service-discovery"
var USERS_GROUP = "/users"
var OAUTH_CLIENTS_GROUP = "/oauth_clients"
var ANALYTICS_GROUP = "/analytics"
var ELASTICSEARCH_URL = "http://localhost:9200"
var ELASTICSEARCH_LOGS_INDEX = "gapi-logs"

type UrlsConstants struct {
	SERVICE_DISCOVERY_GROUP string `json:SERVICE_DISCOVERY_GROUP`
	ANALYTICS_GROUP         string `json:ANALYTICS_GROUP`
}

/*
func LoadURLConstants() {
	urlsJSON, err := ioutil.ReadFile(CONFIGS_LOCATION + URLS_CONFIG_FILE)

	if err != nil {
		return
	}

	var sc UrlsConstants
	json.Unmarshal(urlsJSON, &sc)

	SERVICE_DISCOVERY_URL = "http://" + os.Getenv("SERVICEDISCOVERY_HOST") + ":" + os.Getenv("SERVICEDISCOVERY_PORT")
	SERVICE_DISCOVERY_GROUP = sc.SERVICE_DISCOVERY_GROUP
	ANALYTICS_GROUP = sc.ANALYTICS_GROUP
	ELASTICSEARCH_URL = "http://" + os.Getenv("ELASTICSEARCH_HOST") + ":" + os.Getenv("ELASTICSEARCH_PORT")

	return
}
*/

func LoadURLConstants() {
	SERVICE_DISCOVERY_URL = os.Getenv("API_PROTOCOL") + "://" + os.Getenv("SERVICEDISCOVERY_HOST") + ":" + os.Getenv("SERVICEDISCOVERY_PORT")
	//SERVICE_DISCOVERY_GROUP = GApiConfiguration.Urls.SERVICE_DISCOVERY_GROUP
	//ANALYTICS_GROUP = GApiConfiguration.Urls.ANALYTICS_GROUP
	ELASTICSEARCH_URL = "http://" + os.Getenv("ELASTICSEARCH_HOST") + ":" + os.Getenv("ELASTICSEARCH_PORT")
	if os.Getenv("ELASTICSEARCH_LOGS_INDEX") != "" {
		ELASTICSEARCH_LOGS_INDEX = os.Getenv("ELASTICSEARCH_LOGS_INDEX")
	}
}
