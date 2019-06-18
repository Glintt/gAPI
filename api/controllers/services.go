package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	sdUtils "github.com/Glintt/gAPI/api/servicediscovery/utils"
	"github.com/Glintt/gAPI/api/users"
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
func InternalServiceDiscovery() *servicediscovery.ServiceDiscovery {
	return servicediscovery.GetInternalServiceDiscoveryObject()
}

func ServiceNotFound(c *routing.Context) error {
	http.Response(c, `{"error":true, "msg": "Not found."}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

/*
Normalize all services
*/
func NormalizeServices(c *routing.Context) error {
	user := authentication.GetAuthenticatedUser(c)
	err := ServiceDiscovery(user).NormalizeServices()
	if err != nil {
		http.Response(c, `{"error":true, "msg": "Normalization failed."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return err
	}
	http.Response(c, `{"error":false, "msg": "Normalization done."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateHandler(c *routing.Context) error {
	user := authentication.GetAuthenticatedUser(c)
	serviceID := c.Param("service_id")

	s, err := servicediscovery.ValidateServiceBody(c)
	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	s.Id = bson.ObjectIdHex(serviceID)

	_, err = servicediscovery.ValidateServiceExists(s)
	if err != nil {
		http.Response(c, string(err.Error()), 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	s, err = ServiceDiscovery(user).UpdateService(s)
	if err != nil {
		http.Response(c, `{"error" : true, "msg":"`+err.Error()+`"`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Service updated successfuly"}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func AutoRegisterHandler(c *routing.Context) error {
	serviceDiscovery := InternalServiceDiscovery()
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
	if err != nil {
		// TODO: create service on service discovery
		serviceFound, err = serviceDiscovery.CreateService(serv)
	} else {
		serviceFound.Hosts = append(serviceFound.Hosts, host)
		serviceFound, err = serviceDiscovery.UpdateService(serviceFound)
	}

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "`+err.Error()+`"}`, 400, "AUTO_REGISTER", "application/json")
		return nil
	}

	serv, _ = serviceDiscovery.FindService(serviceFound)
	s2, _ := json.Marshal(serv)

	http.Response(c, string(s2), 201, "AUTO_REGISTER", "application/json")
	return nil
}

func AutoDeRegisterHandler(c *routing.Context) error {
	serviceDiscovery := InternalServiceDiscovery()
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
		serviceFound, err = serviceDiscovery.UpdateService(serviceFound)
	}
	serv, _ = serviceDiscovery.FindService(serviceFound)
	s2, _ := json.Marshal(serv)

	http.Response(c, string(s2), status, "AUTO_DEREGISTER", "application/json")
	return nil
}

func RegisterHandler(c *routing.Context) error {
	user := authentication.GetAuthenticatedUser(c)
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
	serv, err = ServiceDiscovery(user).CreateService(serv)

	resp := `{"error": true, "msg": "Service registered successfuly"}`
	status := 201
	if err != nil {
		resp = `{"error": true, "msg": "` + err.Error() + `"}`
		status = 400
	}

	http.Response(c, resp, status, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func ListServicesHandler(c *routing.Context) error {
	user := authentication.GetAuthenticatedUser(c)
	serviceDiscovery := ServiceDiscovery(user)
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
	services := serviceDiscovery.ListServices(page, searchQuery)

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
	serviceDiscovery := ServiceDiscovery(user)
	matchingURI := c.QueryArgs().Peek("uri")

	// If identifier is passed, search by identifier instead
	identifier := string(c.QueryArgs().Peek("identifier"))
	if identifier != "" {
		serv := service.Service{
			Identifier: string(identifier),
		}
		var err error
		serv, err = serviceDiscovery.FindService(serv)
		if err == nil {
			serviceJSON, _ := json.Marshal(serv)
			http.Response(c, string(serviceJSON), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
			return nil
		}
		http.Response(c, `{"error": true, "msg": "Not found"}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	service, err := serviceDiscovery.GetEndpointForUri(string(matchingURI))

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
	user := authentication.GetAuthenticatedUser(c)
	serviceID := c.Param("service_id")
	serviceDiscovery := ServiceDiscovery(user)
	
	s := service.Service{Id: bson.ObjectIdHex(serviceID)}

	err := serviceDiscovery.DeleteService(s)

	resp := `{"error": true, "msg": "Service deleted successfuly"}`
	status := 200
	if err != nil {
		resp = `{"error": true, "msg": "` + err.Error() + `"}`
		status = 404
	}
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
