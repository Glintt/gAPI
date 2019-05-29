package service

import (
	"errors"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/servicediscovery/servicegroup"
	sdUtils "github.com/Glintt/gAPI/api/servicediscovery/utils"
	"github.com/Glintt/gAPI/api/utils"
	"math/rand"
	"net"

	"gopkg.in/mgo.v2/bson"

	"github.com/Glintt/gAPI/api/http"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

type Service struct {
	Id                         bson.ObjectId `bson:"_id" json:"Id"`
	Identifier                 string
	Name                       string
	Hosts                      []string
	Domain                     string
	Port                       string
	MatchingURI                string
	MatchingURIRegex           string
	ToURI                      string
	Protected                  bool
	APIDocumentation           string
	IsCachingActive            bool
	IsActive                   bool
	HealthcheckUrl             string
	LastActiveTime             int64
	ServiceManagementHost      string
	ServiceManagementPort      string
	ServiceManagementEndpoints map[string]string
	RateLimit                  int
	RateLimitExpirationTime    int64
	IsReachable                bool
	GroupId                    bson.ObjectId `bson:"groupid,omitempty" json:"GroupId"`
	GroupVisibility            bool
	UseGroupAttributes         bool
	ProtectedExclude           map[string]string
}

func Contains(array []int, value int) bool {
	for _, v := range array {
		if v == value {
			return true
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
	if service.Hosts == nil || len(service.Hosts) == 0 {
		return service.Domain + ":" + service.Port
	}

	return service.BalanceUrl()
}

func (service *Service) Call(method string, uri string, headers map[string]string, body string) *fasthttp.Response {
	uri = strings.Replace(uri, service.MatchingURI, service.ToURI, 1)

	callURLWithoutProtocol := service.GetHost() + uri
	callURLWithoutProtocol = strings.Replace(callURLWithoutProtocol, "//", "/", -1)

	callURL := "http://" + callURLWithoutProtocol

	return http.MakeRequest(method, callURL, body, headers)
}

func (service *Service) GenerateId() bson.ObjectId {
	return bson.NewObjectId()
}

func (service *Service) GenerateIdentifier() string {
	identifier := strings.ToLower(service.Name)
	identifier = strings.Replace(identifier, " ", "-", -1)
	return identifier
}

func (service *Service) NormalizeService() {
	service.MatchingURIRegex = sdUtils.GetMatchingURIRegex(service.MatchingURI)
	if service.Id == "" {
		service.Id = service.GenerateId()
	}
	if service.Identifier == "" {
		service.Identifier = service.GenerateIdentifier()
	}
}

func (service *Service) ServiceManagementCall(managementType string) (bool, string) {
	method := service.GetManagementEndpointMethod(managementType)
	callURL := "http://" + service.ServiceManagementHost + ":" + service.ServiceManagementPort + service.GetManagementEndpoint(managementType)

	if ValidateURL(callURL) {
		resp := http.MakeRequest(method, callURL, "", nil)
		utils.LogMessage(string(resp.Body()), utils.DebugLogType)
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

func (service *Service) GetGroup() (servicegroup.ServiceGroup, error) {

	groupId := service.GroupId.Hex()

	servicesGroup, err := servicegroup.ServiceGroupMethods[constants.SD_TYPE]["getbyid"].(func(string) (servicegroup.ServiceGroup, error))(groupId)

	if err != nil {
		return servicegroup.ServiceGroup{}, err
	}

	return servicesGroup, nil
}

func FindServiceInList(s Service, services []Service) (Service, error) {
	for _, rs := range services {
		if rs.MatchingURIRegex == "" {
			rs.MatchingURIRegex = sdUtils.GetMatchingURIRegex(rs.MatchingURI)
		}

		rs.Identifier = rs.GenerateIdentifier()
		// s.Identifier = s.GenerateIdentifier()

		// ser, _ := json.Marshal(rs)

		re := regexp.MustCompile(rs.MatchingURIRegex)
		if re.MatchString(s.MatchingURI) || rs.Id.Hex() == s.Id.Hex() || rs.Identifier == s.Identifier {
			return rs, nil
		}
	}

	return s, errors.New("Not found.")
}
