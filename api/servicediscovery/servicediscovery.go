package servicediscovery

import (
	"encoding/json"
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

var SD_TYPE = "file"

var funcMap = map[string]map[string]interface{}{
	"mongo": {
		"delete": DeleteServiceMongo,
		"update": UpdateMongo,
		"create": CreateServiceMongo,
		"list":   ListServicesMongo,
		"get":    FindMongo},
	"file": {
		"delete": DeleteServiceFile,
		"update": UpdateFile,
		"create": CreateServiceFile,
		"list":   ListServicesFile,
		"get":    FindFile}}

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
	sd.sdAPI.Post("/update", authentication.AuthorizationMiddleware, UpdateHandler)
	sd.sdAPI.Get("/services", ListServicesHandler)
	sd.sdAPI.Get("/endpoint", GetEndpointHandler)
	sd.sdAPI.Delete("/delete", authentication.AuthorizationMiddleware, DeleteEndpointHandler)
	sd.isService = true
}

func UpdateHandler(c *routing.Context) error {
	service, err := ValidateServiceBody(c)

	if err != nil {
		http.Response(c, err.Error(), 400, SERVICE_NAME)
		return nil
	}

	serviceExists, err1 := ValidateServiceExists(service)

	if err1 != nil {
		http.Response(c, string(err1.Error()), 404, SERVICE_NAME)
		return nil
	}

	resp, status := funcMap[SD_TYPE]["update"].(func(Service, Service) (string, int))(service, serviceExists)

	http.Response(c, resp, status, SERVICE_NAME)
	return nil
}

func RegisterHandler(c *routing.Context) error {
	service, err := ValidateServiceBody(c)

	if err != nil {
		http.Response(c, err.Error(), 400, SERVICE_NAME)
		return nil
	}

	_, err = ValidateServiceExists(service)

	// if service exists, return error
	if err == nil {
		http.Response(c, `{"error":true, "msg": "Service already exists."}`, 400, SERVICE_NAME)
		return nil
	}

	resp, status := funcMap[SD_TYPE]["create"].(func(Service) (string, int))(service)

	http.Response(c, resp, status, SERVICE_NAME)
	return nil
}

func ListServicesHandler(c *routing.Context) error {
	services := funcMap[SD_TYPE]["list"].(func() []Service)()

	if len(services) == 0 {
		http.Response(c, `[]`, 200, SERVICE_NAME)
		c.Response.Header.SetContentType("application/json")
		return nil
	}

	list, err := json.Marshal(services)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error parsing body."}`, 404, SERVICE_NAME)
		return nil
	}

	http.Response(c, string(list), 200, SERVICE_NAME)
	c.Response.Header.SetContentType("application/json")
	return nil
}

func GetEndpointHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("uri")

	/*
		fmt.Println("\n=============================================================")
		fmt.Println("SERVICE DISCOVERY =====> uri param = " + string(matchingURI))
		fmt.Println("=============================================================\n")
	*/

	service, err := sd.GetEndpointForUri(string(matchingURI))
	serviceJSON, err1 := json.Marshal(service)

	if err == nil && err1 == nil {
		http.Response(c, string(serviceJSON), 200, SERVICE_NAME)
		return nil
	}
	http.Response(c, `{"error": true, "msg": "Not found"}`, 404, SERVICE_NAME)
	return nil

}

func DeleteEndpointHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("uri")

	resp, status := funcMap[SD_TYPE]["delete"].(func(string) (string, int))(string(matchingURI))

	http.Response(c, resp, status, SERVICE_NAME)
	return nil
}

func (service *ServiceDiscovery) SetIsService(isServ bool) {
	service.isService = isServ
}

func (service *ServiceDiscovery) IsService() bool {
	return service.isService
}
