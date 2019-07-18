package service

import (
	"errors"
	"math/rand"
	"net"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/servicediscovery/servicegroup"
	sdUtils "github.com/Glintt/gAPI/api/servicediscovery/utils"
	userModels "github.com/Glintt/gAPI/api/users/models"
	"github.com/Glintt/gAPI/api/utils"

	"gopkg.in/mgo.v2/bson"

	"regexp"
	"strings"

	"github.com/Glintt/gAPI/api/http"

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
	OAuthClientsEnabled        bool
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
		dialAddress := strings.Replace(host, "http://", "", -1)
		dialAddress = strings.Replace(dialAddress, "https://", "", -1)
		_, err := net.Dial("tcp", dialAddress)
		if err == nil {
			return host
		}
		utils.LogMessage("Dial error: "+err.Error(), utils.DebugLogType)
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

	// callURLWithoutProtocol = strings.Replace(callURLWithoutProtocol, "//", "/", -1)
	callURL := callURLWithoutProtocol
	// If it doesn't have protocol, add http
	if !strings.Contains(callURLWithoutProtocol, "http://") && !strings.Contains(callURLWithoutProtocol, "https://") {
		callURL = "http://" + callURL
	}

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
	groupID := service.GroupId.Hex()
	// Create service group service
	serviceGroupService, err := servicegroup.NewServiceGroupServiceWithUser(userModels.GetInternalAPIUser())
	if err != nil {
		return servicegroup.ServiceGroup{}, err
	}

	// get service group by id
	servicesGroup, err := serviceGroupService.GetServiceGroupById(groupID)
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
		matchinguriWithoutQueryParams := strings.Split(s.MatchingURI, "?")
		matchinguriDecoded := matchinguriWithoutQueryParams[0]

		re := regexp.MustCompile(rs.MatchingURIRegex)
		if re.MatchString(matchinguriDecoded) || rs.Id.Hex() == s.Id.Hex() || rs.Identifier == s.Identifier {
			return rs, nil
		}
	}

	return s, errors.New("Not found.")
}
