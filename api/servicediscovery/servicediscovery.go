package servicediscovery

import (
	"strconv"
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
var PAGE_LENGTH = 10
var SD_TYPE = "file"

var Methods = map[string]map[string]interface{}{
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
	sd.sdAPI.Post("/services/manage", ManageServiceHandler)
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

	resp, status := Methods[SD_TYPE]["update"].(func(Service, Service) (string, int))(service, serviceExists)

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

	resp, status := Methods[SD_TYPE]["create"].(func(Service) (string, int))(service)

	http.Response(c, resp, status, SERVICE_NAME)
	return nil
}

func ListServicesHandler(c *routing.Context) error {
	page := 1
	searchQuery := ""
	if c.QueryArgs().Has("page") {
		var err error
		page, err = strconv.Atoi(string(c.QueryArgs().Peek("page")))

		if err != nil {
			http.Response(c, `{"error" : true, "msg": "Invalid page provided."}`, 404, SERVICE_NAME)
			return nil
		}
	}
	if c.QueryArgs().Has("q") {
		searchQuery = string(c.QueryArgs().Peek("q"))
	}
	
	services := Methods[SD_TYPE]["list"].(func(int, string) []Service)(page, searchQuery)

	if len(services) == 0 {
		http.Response(c, `[]`, 200, SERVICE_NAME)
		c.Response.Header.SetContentType("application/json")
		return nil
	}

	list, jsonErr := json.Marshal(services)

	if jsonErr != nil {
		http.Response(c, `{"error" : true, "msg": "Error parsing body."}`, 404, SERVICE_NAME)
		return nil
	}

	http.Response(c, string(list), 200, SERVICE_NAME)
	c.Response.Header.SetContentType("application/json")
	return nil
}

func GetEndpointHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("uri")

	service, err := sd.GetEndpointForUri(string(matchingURI))
	serviceJSON, err1 := json.Marshal(service)

	if err == nil && err1 == nil {
		http.Response(c, string(serviceJSON), 200, SERVICE_NAME)
		return nil
	}
	http.Response(c, `{"error": true, "msg": "Not found"}`, 404, SERVICE_NAME)
	return nil
}

func ManageServiceHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("service")
	managementType := string(c.QueryArgs().Peek("action"))

	service, err := sd.GetEndpointForUri(string(matchingURI))

	if err == nil {
		success, callResponse := service.ServiceManagementCall(managementType) 
		if success {
			http.Response(c, `{"error": false, "msg": "Service ` + managementType + ` successfuly.", "service_response": "` + callResponse + `"}` , 200, SERVICE_NAME)
			return nil
		}
		http.Response(c, `{"error": true, "msg": "Service could not be ` + managementType + `."}`, 400, SERVICE_NAME)
		return nil
	}
	http.Response(c, `{"error": true, "msg": "Not found"}`, 404, SERVICE_NAME)
	return nil
}

func DeleteEndpointHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("uri")

	resp, status := Methods[SD_TYPE]["delete"].(func(string) (string, int))(string(matchingURI))

	http.Response(c, resp, status, SERVICE_NAME)
	return nil
}

func (service *ServiceDiscovery) SetIsService(isServ bool) {
	service.isService = isServ
}

func (service *ServiceDiscovery) IsService() bool {
	return service.isService
}
