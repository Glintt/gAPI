package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/users"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	sdUtils "github.com/Glintt/gAPI/api/servicediscovery/utils"
	"github.com/Glintt/gAPI/api/utils"

	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

func ServiceDiscoveryServiceName() string {
	return constants.SERVICE_NAME
}

func ServiceDiscovery(user users.User) *servicediscovery.ServiceDiscovery {
	return servicediscovery.GetServiceDiscoveryObject(user)
}

func ServiceNotFound(c *routing.Context) error {
	http.Response(c, `{"error":true, "msg": "Not found."}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func NormalizeServices(c *routing.Context) error {
	err := service.GetServicesRepository(users.User{}).NormalizeServices()
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

	resp, status := service.GetServicesRepository(users.User{}).Update(s, serviceExists)

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
		Hosts:            []string{host},
		MatchingURI:      s["MatchingUri"],
		ToURI:            s["ToUri"],
		Name:             s["Name"],
		HealthcheckUrl:   "http://" + host + s["HealthcheckUrl"],
		APIDocumentation: "http://" + host + s["APIDocumentation"],
		IsReachable:      true,
	}
	serv.MatchingURIRegex = sdUtils.GetMatchingURIRegex(serv.MatchingURI)

	serviceFound, err := servicediscovery.ValidateServiceExists(serv)
	var status int
	var msg string
	if err != nil {
		msg, status = service.GetServicesRepository(users.User{}).CreateService(serv)
	} else {
		serviceFound.Hosts = append(serviceFound.Hosts, host)
		msg, status = service.GetServicesRepository(users.User{}).Update(serviceFound, serviceFound)
	}

	if status > 300 {
		http.Response(c, string(msg), status, "AUTO_REGISTER", "application/json")
		return nil
	}

	serv, _ = service.GetServicesRepository(users.User{}).Find(serv)
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
		_, status = service.GetServicesRepository(users.User{}).Update(serviceFound, serviceFound)
	}
	serv, _ = service.GetServicesRepository(users.User{}).Find(serv)
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
	resp, status := service.GetServicesRepository(users.User{}).CreateService(serv)

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func ListServicesHandler(c *routing.Context) error {
	user := authentication.GetAuthenticatedUser(c)
	var err error
	page := 1
	searchQuery := ""

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
	
	services := service.GetServicesRepository(user).ListServices(page, searchQuery)

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
	user := authentication.GetAuthenticatedUser(c)
	matchingURI := c.QueryArgs().Peek("uri")

	// If identifier is passed, search by identifier instead
	identifier := string(c.QueryArgs().Peek("identifier"))
	if identifier != "" {
		serv := service.Service{
			Identifier: string(identifier),
		}
		var err error
		serv, err = service.GetServicesRepository(user).Find(serv)
		if err == nil {
			serviceJSON, _ := json.Marshal(serv)
			http.Response(c, string(serviceJSON), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
			return nil
		}
		http.Response(c, `{"error": true, "msg": "Not found"}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	utils.LogMessage("HERE", utils.DebugLogType)

	service, err := ServiceDiscovery(user).GetEndpointForUri(string(matchingURI))

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

	resp, status := service.GetServicesRepository(users.User{}).DeleteService(s)

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func ManageServiceHandler(c *routing.Context) error {
	user := authentication.GetAuthenticatedUser(c)
	matchingURI := c.QueryArgs().Peek("service")
	managementType := string(c.QueryArgs().Peek("action"))

	service, err := ServiceDiscovery(user).GetEndpointForUri(string(matchingURI))

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