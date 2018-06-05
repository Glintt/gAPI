package servicediscovery

import (
	"fmt"
	"gAPIManagement/api/authentication"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"github.com/qiangxue/fasthttp-routing"
)

type ServiceDiscovery struct {
	isService          bool
	registeredServices []Service
	sdAPI              *routing.RouteGroup
}

var sd ServiceDiscovery

var SERVICE_NAME = "/service-discovery"
var PAGE_LENGTH = 10
var SD_TYPE = "file"

var Methods = map[string]map[string]interface{}{
	"mongo": {
		"delete": DeleteServiceMongo,
		"update": UpdateMongo,
		"create": CreateServiceMongo,
		"list":   ListServicesMongo,
		"get":    FindMongo,
		"normalize": NormalizeServicesMongo},
	"file": {
		"delete": DeleteServiceFile,
		"update": UpdateFile,
		"create": CreateServiceFile,
		"list":   ListServicesFile,
		"get":    FindFile,
		"normalize": NormalizeServicesFile}}

func (serviceDisc *ServiceDiscovery) SetRegisteredServices(rs []Service) {
	serviceDisc.registeredServices = rs
}

func GetServiceDiscoveryObject() *ServiceDiscovery {
	return &sd
}

func StartServiceDiscovery(router *routing.Router) {
	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" {
		SD_TYPE = "mongo"
		InitMongo()
	} else {
		servicesConfig := LoadServicesConfiguration()
		sd.registeredServices = servicesConfig.Services
	}

	sd.sdAPI = router.Group(config.SERVICE_DISCOVERY_GROUP)

	sd.isService = true
	sd.sdAPI.Post("/register", authentication.AuthorizationMiddleware, RegisterHandler)
	sd.sdAPI.Post("/admin/normalize", authentication.AuthorizationMiddleware, NormalizeServices)
	sd.sdAPI.Post("/update", authentication.AuthorizationMiddleware, UpdateHandler)
	sd.sdAPI.Get("/services", ListServicesHandler)
	sd.sdAPI.Get("/endpoint", GetEndpointHandler)
	sd.sdAPI.Delete("/delete", authentication.AuthorizationMiddleware, DeleteEndpointHandler)
	sd.sdAPI.Post("/services/manage", ManageServiceHandler)
	sd.sdAPI.Get("/services/manage/types", ManageServiceTypesHandler)
	sd.sdAPI.Post("/service-groups/register", authentication.AuthorizationMiddleware, RegisterServiceGroupHandler)
	// sd.sdAPI.Post("/service-groups/service/register", authentication.AuthorizationMiddleware, RegisterServiceToServiceGroupHandler)
	sd.sdAPI.Get("/service-groups", ListServiceGroupsHandler)
	sd.sdAPI.To("GET,POST,PUT,PATCH,DELETE", "/*", ServiceNotFound)
	sd.isService = true
}


func DeleteEndpointHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("uri")

	service := Service{MatchingURI: string(matchingURI)}
	resp, status := Methods[SD_TYPE]["delete"].(func(Service) (string, int))(service)

	http.Response(c, resp, status, SERVICE_NAME)
	return nil
}

func (service *ServiceDiscovery) SetIsService(isServ bool) {
	service.isService = isServ
}

func (service *ServiceDiscovery) IsService() bool {
	return service.isService
}


func (service *ServiceDiscovery) IsExternalRequest(requestContxt *routing.Context) bool {
	hosts, _ := ListAllAvailableHosts()
	requestHost := string(requestContxt.Host())
	fmt.Println(hosts)

	for _, v := range hosts {
		if v == requestHost {
			return false
		}
	}
	return true
}