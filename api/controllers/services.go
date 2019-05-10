package controllers

import (
	"encoding/json"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/servicediscovery/service"
	sdUtils "gAPIManagement/api/servicediscovery/utils"
	"gAPIManagement/api/utils"
	"strconv"

	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
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
	serviceID := c.Param("service_id")
	s, err := servicediscovery.ValidateServiceBody(c)

	s.Id = bson.ObjectIdHex(serviceID)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	serviceExists, err1 := servicediscovery.ValidateServiceExists(s)

	if err1 != nil {
		http.Response(c, string(err1.Error()), 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	resp, status := Methods()["update"].(func(service.Service, service.Service) (string, int))(s, serviceExists)

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func AutoRegisterHandler(c *routing.Context) error {
	var s map[string]string
	json.Unmarshal(c.Request.Body(), &s)

	if s["MatchingUri"] == "" || s["ToUri"] == "" || s["Name"] == "" || s["Port"] == "" {
		http.Response(c, `{"error" : true, "msg": "Missing body parameters."}`,
			400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	host := c.RemoteIP().String() + ":" + s["Port"]
	serv := service.Service{
		Hosts:       []string{host},
		MatchingURI: s["MatchingUri"],
		ToURI:       s["ToUri"],
		Name:        s["Name"],
	}
	serv.MatchingURIRegex = sdUtils.GetMatchingURIRegex(serv.MatchingURI)

	serviceFound, err := servicediscovery.ValidateServiceExists(serv)
	var status int
	if err != nil {
		_, status = service.CreateServiceMongo(serv)

	} else {
		serviceFound.Hosts = append(serviceFound.Hosts, host)
		_, status = service.UpdateMongo(serviceFound, serviceFound)
	}
	serv, _ = service.FindMongo(serv)
	s2, _ := json.Marshal(serv)

	http.Response(c, string(s2), status, "AUTO_REGISTER", "application/json")
	return nil
}

func AutoDeRegisterHandler(c *routing.Context) error {
	var s map[string]string
	json.Unmarshal(c.Request.Body(), &s)

	if s["MatchingUri"] == "" || s["Port"] == "" {
		http.Response(c, `{"error" : true, "msg": "Missing body parameters."}`,
			400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	host := c.RemoteIP().String() + ":" + s["Port"]

	serv := service.Service{
		MatchingURI: s["MatchingUri"],
	}
	serv.MatchingURIRegex = sdUtils.GetMatchingURIRegex(serv.MatchingURI)

	serviceFound, err := servicediscovery.ValidateServiceExists(serv)
	var status int
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Service not found."}`,
			400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	} else {
		serviceFound.Hosts = utils.RemoveStringFromArray(serviceFound.Hosts, host)
		_, status = service.UpdateMongo(serviceFound, serviceFound)
	}
	serv, _ = service.FindMongo(serv)
	s2, _ := json.Marshal(serv)

	http.Response(c, string(s2), status, "AUTO_DEREGISTER", "application/json")
	return nil
}

func RegisterHandler(c *routing.Context) error {
	serv, err := servicediscovery.ValidateServiceBody(c)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	_, err = servicediscovery.ValidateServiceExists(serv)

	// if service exists, return error
	if err == nil {
		http.Response(c, `{"error":true, "msg": "Service already exists."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	serv.MatchingURIRegex = sdUtils.GetMatchingURIRegex(serv.MatchingURI)
	resp, status := Methods()["create"].(func(service.Service) (string, int))(serv)

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func ListServicesHandler(c *routing.Context) error {
	var err error
	page := 1
	searchQuery := ""
	user := string(c.Request.Header.Peek("User"))

	if c.QueryArgs().Has("page") {
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

	services := Methods()["list"].(func(int, string, bool) []service.Service)(page, searchQuery, permissions)

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

	// If identifier is passed, search by identifier instead
	identifier := string(c.QueryArgs().Peek("identifier"))
	if identifier != "" {
		serv := service.Service{
			Identifier: string(identifier),
		}
		var err error
		serv, err = Methods()["get"].(func(service.Service) (service.Service, error))(serv)
		if err == nil {
			serviceJSON, _ := json.Marshal(serv)
			http.Response(c, string(serviceJSON), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
			return nil
		}
	}

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
	serviceID := c.Param("service_id")
	// matchingURI := c.QueryArgs().Peek("uri")

	// service := servicediscovery.Service{MatchingURI: string(matchingURI)}
	s := service.Service{Id: bson.ObjectIdHex(serviceID)}

	resp, status := Methods()["delete"].(func(service.Service) (string, int))(s)

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
