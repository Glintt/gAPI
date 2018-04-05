package config


import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var URLS_CONFIG_FILE = "urls.json"
var SERVICE_DISCOVERY_URL = "http://localhost:8080"
var SERVICE_DISCOVERY_GROUP = "/service-discovery"
var ANALYTICS_GROUP = "/analytics"
var ELASTICSEARCH_URL = "http://localhost:9200"

type UrlsConstants struct {
	SERVICE_DISCOVERY_GROUP string `json:SERVICE_DISCOVERY_GROUP`
	ANALYTICS_GROUP         string `json:ANALYTICS_GROUP`
}

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
