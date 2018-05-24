package config

import (
	"encoding/json"
	"io/ioutil"
)

var GAPI_CONFIG_FILE = "gAPI.json"

type GApiConfig struct {
	Authentication   GApiAuthenticationConfig
	Logs             GApiLogsConfig
	Cors             GApiCorsConfig
	ServiceDiscovery GApiServiceDiscoveryConfig
	Urls             UrlsConstants
	Healthcheck      GApiHealthCheckConfig
	Notifications 	 GApiNotificationsConfig
	ManagementTypes  map[string]map[string]string
	RateLimiting 	GApiRateLimitingConfig
}

type GApiAuthenticationConfig struct {
	Username string
	Password string
	TokenExpirationTime int
	TokenSigningKey string
}
type GApiLogsConfig struct {
	Active bool
	Type   string
}

type GApiCorsConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

type GApiServiceDiscoveryConfig struct {
	Type string
}

type GApiHealthCheckConfig struct {
	Active    bool
	Frequency int
	Notification bool
}

type GApiNotificationsConfig struct {
	Type    string
	Slack GApiSlackNotificationsConfig
}

type GApiSlackNotificationsConfig struct {
	WebhookUrl string
}

type GApiRateLimitingConfig struct {
	Active bool
	Limit int
	Period int64
	Metrics []string
}

var GApiConfiguration GApiConfig

func LoadGApiConfig() {
	gapiJSON, err := ioutil.ReadFile(CONFIGS_LOCATION + GAPI_CONFIG_FILE)

	if err != nil {
		return
	}

	json.Unmarshal(gapiJSON, &GApiConfiguration)
	return
}
