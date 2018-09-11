package thirdpartyauthentication

import (
	"errors"
	"gAPIManagement/api/utils"

	// "gAPIManagement/api/config"
	// "gAPIManagement/api/config"
	"gAPIManagement/api/http"

	"github.com/valyala/fasthttp"
)

type OAuthServer struct {
	Host                 string
	Port                 string
	AuthorizeEndpoint    string
	UserTokenInformation UserTokenInformation
}

func GetAuthorizeEndpointUrl(oas OAuthServer) string {
	return oas.Host + ":" + oas.Port + oas.AuthorizeEndpoint
}

/*
func LoadFromConfig() OAuthServer {
	authenticationJSON, err := utils.LoadJsonFile(config.CONFIGS_LOCATION + config.AUTHENTICATION_CONFIG_FILE)

	if err != nil {
		return OAuthServer{}
	}

	var oas OAuthServer
	json.Unmarshal(authenticationJSON, &oas)

	oas.AuthorizeEndpoint = oas.Host + ":" + oas.Port + oas.AuthorizeEndpoint
	return oas
} */

func (oauthServer *OAuthServer) Authorize(request fasthttp.Request) (string, string, error) {

	token := request.Header.Peek("Authorization")

	if token == nil {
		return "", "", errors.New("Not Authorized.")
	}

	// TODO: send on header
	var headers map[string]string
	headers = make(map[string]string)
	headers["Authorization"] = string(token)

	utils.LogMessage("ThirdPartyOAuth (AuthorizationEndpoint) : "+oauthServer.AuthorizeEndpoint, utils.DebugLogType)

	response := http.MakeRequest("GET", oauthServer.AuthorizeEndpoint, "", headers)

	if response.StatusCode() != 200 {
		return "", "", errors.New("Not Authorized.")
	}

	if oauthServer.UserTokenInformation.Active {
		return oauthServer.UserTokenInformation.Name, string(response.Header.Peek(oauthServer.UserTokenInformation.Name)), nil
	}

	return "token", "valid", nil
}
