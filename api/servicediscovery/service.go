package servicediscovery

import (
	//"gAPIManagement/api/servicediscovery"
	"gopkg.in/mgo.v2/bson"
	"net"
	"math/rand"
	"fmt"
	"gAPIManagement/api/config"
	"gAPIManagement/api/utils"
	
	"gAPIManagement/api/http"
	"strings"
	"regexp"
	"github.com/valyala/fasthttp"
	//"gopkg.in/mgo.v2/bson"
)

type Service struct {
	Id                    bson.ObjectId `bson:"_id" json:"Id"`
	Name                  string
	Hosts 				  []string
	Domain				  string
	Port				  string
	MatchingURI           string
	MatchingURIRegex      string
	ToURI                 string
	Protected             bool
	APIDocumentation      string
	IsCachingActive       bool 
	IsActive              bool 
	HealthcheckUrl		  string 
	LastActiveTime 		  int64 
	ServiceManagementHost       string 
	ServiceManagementPort       string 
	ServiceManagementEndpoints  map[string]string
	RateLimit	int
	RateLimitExpirationTime	int64
	IsReachable              bool
}

func Contains(array []int, value int) bool {
	for _, v := range array{
		if v == value{
			return true
		}
	}
	return false
}

func (service *Service) IsReachableFromExternal(sd ServiceDiscovery) bool {
	if service.IsReachable {
		return service.IsReachable
	}

	sgList, err := sd.GetListOfServicesGroup()
	if err != nil {
		return false
	}

	for _, sg := range sgList {
		if sg.Contains(*service) {
			return sg.IsReachable
		}
	}
	return false
}

func (service *Service) BalanceUrl() string {
	numHosts := len(service.Hosts)
	indexesToTry := rand.Perm(numHosts)

	for _, index := range indexesToTry {
		host := service.Hosts[index]
		_, err := net.Dial("tcp", host)
		if err == nil {
			return host
		}
	}

	return service.Domain + ":" + service.Port
}
func (service *Service) GetHost() string {
	if service.Hosts == nil || len(service.Hosts) == 0{
		return service.Domain + ":" + service.Port
	}

	return service.BalanceUrl()
}

func (service *Service) Call(method string, uri string, headers map[string]string, body string) *fasthttp.Response {
	uri = strings.Replace(uri, service.MatchingURI, service.ToURI, 1)

	callURLWithoutProtocol := service.GetHost() + uri
	callURLWithoutProtocol = strings.Replace(callURLWithoutProtocol, "//", "/", -1)
	fmt.Println(callURLWithoutProtocol)
	callURL := "http://" + callURLWithoutProtocol

	return http.MakeRequest(method, callURL, body, headers)
}

func (service *Service) GenerateId() bson.ObjectId {
	return bson.NewObjectId()
}

func (service *Service) NormalizeService() {
	service.MatchingURIRegex = GetMatchingURIRegex(service.MatchingURI)
	if service.Id == "" {
		service.Id = service.GenerateId()
	}
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
	return service.ServiceManagementEndpoints[managementType]
}

func (service *Service) GetManagementEndpointMethod(managementType string) string {
	return config.GApiConfiguration.ManagementTypes[managementType]["method"]
}

func ValidateURL(url string) bool {
	var validURL = regexp.MustCompile(`^(((http|https):\/{2})+(([0-9a-z_-]+\.?)+(:[0-9]+)?((\/([~0-9a-zA-Z#\+%@\.\/_-]+))?(\?[0-9a-zA-Z\+%@\/&\[\];=_-]+)?)?))\b$`)
	return validURL.MatchString(url)
}