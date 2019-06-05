package oauth_clients

import (
	"github.com/Glintt/gAPI/api/config"
)

type OAuthClient struct {
	ClientId     string
	ClientSecret string
}

const (
	OAUTH_CLIENTS_COLLECTION = "gapi_oauth_clients"
	SERVICE_NAME             = "gapi_oauth_clients"
	PAGE_LENGTH              = 10
)

func Find(clientId string, clientSecret string) OAuthClient {
	return OAuthClientsMethods[config.GApiConfiguration.ServiceDiscovery.Type]["find"].(func(string, string) OAuthClient)(clientId, clientSecret)
}
