package servicediscovery

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/utils"
	"encoding/json"
)

type ServicesConfig struct {
	Services []Service `json:"services"`
}

func LoadServicesConfiguration() ServicesConfig {
	servicesJSON, err := utils.LoadJsonFile(config.CONFIGS_LOCATION + config.SERVICE_DISCOVERY_CONFIG_FILE)

	if err != nil {
		return ServicesConfig{}
	}

	var sc ServicesConfig
	json.Unmarshal(servicesJSON, &sc)
	return sc
}
