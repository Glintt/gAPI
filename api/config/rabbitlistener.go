package config

import (
	"encoding/json"
	"io/ioutil"
)

var RABBIT_CONFIG_FILE = "rabbit-listener.json"

type RabbitListenerConfig struct{
	Cors RabbitListenerCorsConfig
}

type RabbitListenerCorsConfig struct {
	AllowedOrigins []string
	AllowCredentials bool
}

var RabbitListenerConfigurationObj RabbitListenerConfig

func LoadRabbitListenerConfig(){
	rlisJSON, err := ioutil.ReadFile(CONFIGS_LOCATION + RABBIT_CONFIG_FILE)

	if err != nil {
		return
	}

	json.Unmarshal(rlisJSON, &RabbitListenerConfigurationObj)
	return
}