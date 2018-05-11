package servicediscovery

import (
	"gAPIManagement/api/utils"
	
	"gAPIManagement/api/http"
	"strings"
	"regexp"
	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

type Service struct {
	ID                    bson.ObjectId `bson:"_id" json:"id"`
	Name                  string
	Domain                string
	Port                  string
	MatchingURI           string
	ToURI                 string
	Protected             bool
	APIDocumentation      string
	CachingExpirationTime int 
	IsCachingActive       bool 
	IsActive              bool 
	HealthcheckUrl	string 
	LastActiveTime int64 
	ServiceManagementHost       string 
	ServiceManagementPort       string 
	RestartEndpoint       string  
	RedeployEndpoint       string 
	UndeployEndpoint       string 
	BackupEndpoint       string
	LogsEndpoint string
}

var ManagementTypes = map[string]map[string]string{
	"restart" : {
		"action": "restart",
		"method": "POST"},
	"undeploy" : {
		"action": "undeploy",
		"method": "POST"},
	"redeploy" : {
		"action": "redeploy",
		"method": "POST"},
	"backup" : {
		"action": "backup",
		"method": "POST"},
	"logs" : {
		"action": "logs",
		"method": "GET"}}

func (service *Service) Call(method string, uri string, headers map[string]string, body string) *fasthttp.Response {
	uri = strings.Replace(uri, service.MatchingURI, service.ToURI, 1)

	callURLWithoutProtocol := service.Domain + ":" + service.Port + uri
	callURLWithoutProtocol = strings.Replace(callURLWithoutProtocol, "//", "/", -1)

	callURL := "http://" + callURLWithoutProtocol

	return http.MakeRequest(method, callURL, body, headers)
}

func (service *Service) ServiceManagementCall(managementType string) (bool, string) {
	method := service.GetManagementEndpointMethod(managementType)
	callURL := "http://" + service.ServiceManagementHost + ":" + service.ServiceManagementPort + service.GetManagementEndpoint(managementType)

	if ValidateURL(callURL) {
		resp := http.MakeRequest(method, callURL, "", nil)
		utils.LogMessage(string(resp.Body()))
		if resp.StatusCode() != 200 {
			return false, string(resp.Body())
		}
		return true, string(resp.Body())
	}
	return false, ""
}

func (service *Service) GetManagementEndpoint(managementType string) string {
	switch managementType {
		case ManagementTypes["restart"]["action"]:
			return service.RestartEndpoint
		case ManagementTypes["undeploy"]["action"]:
			return service.UndeployEndpoint
		case ManagementTypes["redeploy"]["action"]:
			return service.RedeployEndpoint
		case ManagementTypes["backup"]["action"]:
			return service.BackupEndpoint
		case ManagementTypes["logs"]["action"]:
			return service.LogsEndpoint
	}
	return ":"
}

func (service *Service) GetManagementEndpointMethod(managementType string) string {
	return ManagementTypes[managementType]["method"]
}

func ValidateURL(url string) bool {
	var validURL = regexp.MustCompile(`^(((http|https):\/{2})+(([0-9a-z_-]+\.?)+(:[0-9]+)?((\/([~0-9a-zA-Z#\+%@\.\/_-]+))?(\?[0-9a-zA-Z\+%@\/&\[\];=_-]+)?)?))\b$`)
	return validURL.MatchString(url)
}