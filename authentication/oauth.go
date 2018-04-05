package authentication

import (
	"api-management/config"
	"api-management/http"
	"api-management/utils"
	"encoding/json"
	"errors"

	"github.com/valyala/fasthttp"
)

type OAuthServer struct {
	Host                 string               `json:"host"`
	Port                 string               `json:"port"`
	AuthorizeEndpoint    string               `json:"authorize_endpoint"`
	UserTokenInformation UserTokenInformation `json:"token_user_information"`
}

func LoadFromConfig() OAuthServer {
	authenticationJSON, err := utils.LoadJsonFile(config.CONFIGS_LOCATION + config.AUTHENTICATION_CONFIG_FILE)

	if err != nil {
		return OAuthServer{}
	}

	var oas OAuthServer
	json.Unmarshal(authenticationJSON, &oas)

	oas.AuthorizeEndpoint = oas.Host + ":" + oas.Port + oas.AuthorizeEndpoint
	return oas
}

func (oauthServer *OAuthServer) Authorize(request fasthttp.Request) (string, string, error) {

	token := request.Header.Peek("Authorization")

	if token == nil {
		return "", "", errors.New("Not Authorized.")
	}

	// TODO: send on header
	var headers map[string]string
	headers = make(map[string]string)
	headers["Authorization"] = string(token)

	response := http.MakeRequest(config.GET, oauthServer.AuthorizeEndpoint, "", headers)

	if response.StatusCode() != 200 {
		return "", "", errors.New("Not Authorized.")
	}

	if oauthServer.UserTokenInformation.Active {
		return oauthServer.UserTokenInformation.Name, string(response.Header.Peek(oauthServer.UserTokenInformation.Name)), nil
	}

	return "token", "valid", nil
}