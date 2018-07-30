package controllers

import (
	"encoding/json"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery"
	"strconv"

	"github.com/qiangxue/fasthttp-routing"
)

func Methods() map[string]interface{} {
	return servicediscovery.Methods[servicediscovery.SD_TYPE]
}

func ServiceDiscoveryServiceName() string {
	return servicediscovery.SERVICE_NAME
}

func ServiceDiscovery() *servicediscovery.ServiceDiscovery {
	return servicediscovery.GetServiceDiscoveryObject()
}

func ServiceNotFound(c *routing.Context) error {
	http.Response(c, `{"error":true, "msg": "Not found."}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func NormalizeServices(c *routing.Context) error {
	err := Methods()["normalize"].(func() error)()
	if err != nil {
		http.Response(c, `{"error":true, "msg": "Normalization failed."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return err
	}
	http.Response(c, `{"error":false, "msg": "Normalization done."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateHandler(c *routing.Context) error {
	service, err := servicediscovery.ValidateServiceBody(c)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	serviceExists, err1 := servicediscovery.ValidateServiceExists(service)

	if err1 != nil {
		http.Response(c, string(err1.Error()), 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	resp, status := Methods()["update"].(func(servicediscovery.Service, servicediscovery.Service) (string, int))(service, serviceExists)

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func RegisterHandler(c *routing.Context) error {
	service, err := servicediscovery.ValidateServiceBody(c)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	_, err = servicediscovery.ValidateServiceExists(service)

	// if service exists, return error
	if err == nil {
		http.Response(c, `{"error":true, "msg": "Service already exists."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	service.MatchingURIRegex = servicediscovery.GetMatchingURIRegex(service.MatchingURI)
	resp, status := Methods()["create"].(func(servicediscovery.Service) (string, int))(service)

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func ListServicesHandler(c *routing.Context) error {
	page := 1
	searchQuery := ""
	user := string(c.Request.Header.Peek("User"))

	if c.QueryArgs().Has("page") {
		var err error
		page, err = strconv.Atoi(string(c.QueryArgs().Peek("page")))

		if err != nil {
			http.Response(c, `{"error" : true, "msg": "Invalid page provided."}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
			return nil
		}
	}
	if c.QueryArgs().Has("q") {
		searchQuery = string(c.QueryArgs().Peek("q"))
	}

	permissions := false
	if user != "" {
		permissions = true
	}

	services := Methods()["list"].(func(int, string, bool) []servicediscovery.Service)(page, searchQuery, permissions)

	if len(services) == 0 {
		http.Response(c, `[]`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		c.Response.Header.SetContentType("application/json")
		return nil
	}

	list, jsonErr := json.Marshal(services)

	if jsonErr != nil {
		http.Response(c, `{"error" : true, "msg": "Error parsing body."}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, string(list), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	c.Response.Header.SetContentType("application/json")
	return nil
}

func GetEndpointHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("uri")

	service, err := ServiceDiscovery().GetEndpointForUri(string(matchingURI))

	group, getGroupErr := service.GetGroup()
	if getGroupErr != nil {
		service.GroupVisibility = service.IsReachable
	} else {
		service.GroupVisibility = group.IsReachable
	}

	serviceJSON, err1 := json.Marshal(service)

	if err == nil && err1 == nil {
		http.Response(c, string(serviceJSON), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error": true, "msg": "Not found"}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func DeleteEndpointHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("uri")

	service := servicediscovery.Service{MatchingURI: string(matchingURI)}
	resp, status := Methods()["delete"].(func(servicediscovery.Service) (string, int))(service)

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func ManageServiceHandler(c *routing.Context) error {
	matchingURI := c.QueryArgs().Peek("service")
	managementType := string(c.QueryArgs().Peek("action"))

	service, err := ServiceDiscovery().GetEndpointForUri(string(matchingURI))

	if err == nil {
		success, callResponse := service.ServiceManagementCall(managementType)

		if success {
			http.Response(c, `{"error": false, "msg": "Service `+managementType+` successfuly.", "service_response": `+strconv.Quote(callResponse)+`}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
			return nil
		}
		http.Response(c, `{"error": true, "msg": "Service could not be `+managementType+`."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error": true, "msg": "Not found"}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func ManageServiceTypesHandler(c *routing.Context) error {
	managementTypesJson, err := json.Marshal(config.GApiConfiguration.ManagementTypes)

	response := string(managementTypesJson)
	statusCode := 200
	if err != nil {
		response = `{"error": true, "msg": "Not found"}`
		statusCode = 404
	}

	http.Response(c, response, statusCode, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}
