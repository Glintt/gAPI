package config

var CONFIGS_LOCATION = "./configs/"

var SERVICE_DISCOVERY_CONFIG_FILE = "services.json"
var AUTHENTICATION_CONFIG_FILE = "oauth.json"

var POST = "POST"
var GET = "GET"
var DELETE = "DELETE"
var PUT = "PUT"
var PATCH = "PATCH"

const GAPI_API_LOGS_INDEX = "gapi-api-logs"

var APPLICATION_JSON = "application/json"

var SOCKET_PORT_DEFAULT = "5000"

var MATCHING_URI_REGEX = "((/([\\w?\\-=:.&+#])*)*$)"

func LoadConfigs() {
	LoadGApiConfig()
	LoadURLConstants()
}
