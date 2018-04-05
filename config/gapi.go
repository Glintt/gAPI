package config

import (
	"encoding/json"
	"io/ioutil"
)

var GAPI_CONFIG_FILE = "gAPI.json"


type GApiConfig struct {
	Authentication GApiAuthenticationConfig
	Logs GApiLogsConfig
}

type GApiAuthenticationConfig struct {
	Username string
	Password string
}
type GApiLogsConfig struct{
	Active bool
	Type string
}

var GApiConfiguration GApiConfig

func LoadGApiConfig(){
	gapiJSON, err := ioutil.ReadFile(CONFIGS_LOCATION + GAPI_CONFIG_FILE)

	if err != nil {
		return
	}

	json.Unmarshal(gapiJSON, &GApiConfiguration)
	return
}