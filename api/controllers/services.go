package controllers

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	sdUtils "github.com/Glintt/gAPI/api/servicediscovery/utils"
	"github.com/Glintt/gAPI/api/utils"

	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

// ServiceDiscoveryServiceName Name of group of services to use on log storage
func ServiceDiscoveryServiceName() string {
	return constants.SERVICE_NAME
}

// ServiceDiscovery return service discovery object with user context, only with access to user's resources
func ServiceDiscovery(c *routing.Context) *servicediscovery.ServiceDiscovery {
	user := authentication.GetAuthenticatedUser(c)
	return servicediscovery.GetServiceDiscoveryObject(user)
}

// InternalServiceDiscovery return service discovery object as internal, with access to everything
func InternalServiceDiscovery() *servicediscovery.ServiceDiscovery {
	return servicediscovery.GetInternalServiceDiscoveryObject()
}

// ServiceNotFound return Service not found response
func ServiceNotFound(c *routing.Context) error {
	return http.NotFound(c, "Service not found.", ServiceDiscoveryServiceName())
}

// NormalizeServices normalizes all services so theyt match the specified rules
func NormalizeServices(c *routing.Context) error {
	err := ServiceDiscovery(c).NormalizeServices()
	if err != nil {
		return http.Error(c, "Normalization failed", 400, ServiceDiscoveryServiceName())
	}
	return http.OkFormated(c, "Normalization done.", ServiceDiscoveryServiceName())
}

// UpdateHandler updates the service
func UpdateHandler(c *routing.Context) error {
	// Get service discovery object
	serviceDiscovery := ServiceDiscovery(c)

	serviceID := c.Param("service_id")

	s, status, err := servicediscovery.ValidateServiceBodyAndServiceExists(c, serviceID)
	if err != nil {
		return http.Error(c, err.Error(), status, ServiceDiscoveryServiceName())
	}

	s, err = serviceDiscovery.UpdateService(s)
	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	return http.Created(c, "Service updated successfuly", ServiceDiscoveryServiceName())
}

// AutoRegisterHandler handles auto register post request
func AutoRegisterHandler(c *routing.Context) error {
	// Create internal service discovery object with internal api user, which has access to everything
	serviceDiscovery := InternalServiceDiscovery()

	var s map[string]string
	json.Unmarshal(c.Request.Body(), &s)

	if s["MatchingUri"] == "" || s["ToUri"] == "" || s["Name"] == "" || s["Port"] == "" {
		return http.Error(c, "Missing body parameters.", 400, ServiceDiscoveryServiceName())
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
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	serv, _ = serviceDiscovery.FindService(serviceFound)
	s2, _ := json.Marshal(serv)

	http.Response(c, string(s2), 201, "AUTO_REGISTER", "application/json")
	return nil
}

// AutoDeRegisterHandler handles auto deregister post request
func AutoDeRegisterHandler(c *routing.Context) error {
	serviceDiscovery := InternalServiceDiscovery()
	var s map[string]string
	json.Unmarshal(c.Request.Body(), &s)

	if s["MatchingUri"] == "" || s["Port"] == "" {
		return http.Error(c, "Missing body parameters.", 400, ServiceDiscoveryServiceName())
	}
	host := c.RemoteIP().String() + ":" + s["Port"]

	serv := service.Service{
		MatchingURI: s["MatchingUri"],
	}
	serv.MatchingURIRegex = sdUtils.GetMatchingURIRegex(serv.MatchingURI)

	serviceFound, err := servicediscovery.ValidateServiceExists(serv)
	var status int
	if err != nil {
		return ServiceNotFound(c)
	} else {
		serviceFound.Hosts = utils.RemoveStringFromArray(serviceFound.Hosts, host)
		serviceFound, err = serviceDiscovery.UpdateService(serviceFound)
	}
	serv, _ = serviceDiscovery.FindService(serviceFound)
	s2, _ := json.Marshal(serv)

	http.Response(c, string(s2), status, "AUTO_DEREGISTER", "application/json")
	return nil
}

// RegisterHandler handles register post request
func RegisterHandler(c *routing.Context) error {
	// Get service discovery object
	serviceDiscovery := ServiceDiscovery(c)

	// Validate if request body is valid
	serv, err := servicediscovery.ValidateServiceBody(c)
	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	// Validate if service already exists
	_, err = servicediscovery.ValidateServiceExists(serv)
	if err == nil {
		return http.Error(c, `Service already exists.`, 400, ServiceDiscoveryServiceName())
	}

	serv.MatchingURIRegex = sdUtils.GetMatchingURIRegex(serv.MatchingURI)
	serv, err = serviceDiscovery.CreateService(serv)

	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}
	return http.Created(c, "Service registered successfuly", ServiceDiscoveryServiceName())
}

// parseListServicesHandlerParameters parses page and search query parameters for GET /services
func parseListServicesHandlerParameters(c *routing.Context) (string, int, error) {
	page := 1
	var err error
	searchQuery := ""

	if c.QueryArgs().Has("page") {
		page, err = strconv.Atoi(string(c.QueryArgs().Peek("page")))
		if err != nil {
			return "", page, errors.New("Invalid page")
		}
	}
	if c.QueryArgs().Has("q") {
		searchQuery = string(c.QueryArgs().Peek("q"))
	}

	return searchQuery, page, nil
}

// ListServicesHandler handles /services get request
func ListServicesHandler(c *routing.Context) error {
	// create service discovery object with authenticated user
	serviceDiscovery := ServiceDiscovery(c)

	// Get page and search query parameters
	searchQuery, page, err := parseListServicesHandlerParameters(c)
	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	services := serviceDiscovery.ListServices(page, searchQuery)
	if len(services) == 0 {
		return http.Ok(c, `[]`, ServiceDiscoveryServiceName())
	}

	list, err := json.Marshal(services)
	if err != nil {
		return http.Error(c, `Error parsing body.`, 400, ServiceDiscoveryServiceName())
	}

	return http.Ok(c, string(list), ServiceDiscoveryServiceName())
}

// GetEndpointHandler handles GET /endpoint request
func GetEndpointHandler(c *routing.Context) error {
	// create service discovery object with authenticated user
	serviceDiscovery := ServiceDiscovery(c)

	// get query parameters
	matchingURI := string(c.QueryArgs().Peek("uri"))
	identifier := string(c.QueryArgs().Peek("identifier"))

	// Search service by identifier or matchingURI. If error, throw not found response
	serv, err := serviceDiscovery.FindByIdentifierOrMatchingUri(identifier, matchingURI)
	if err != nil {
		return ServiceNotFound(c)
	}

	// get group visibility
	group, getGroupErr := serv.GetGroup()
	if getGroupErr != nil {
		serv.GroupVisibility = serv.IsReachable
	} else {
		serv.GroupVisibility = group.IsReachable
	}

	serviceJSON, err := json.Marshal(serv)
	if err != nil {
		return ServiceNotFound(c)
	}
	return http.Ok(c, string(serviceJSON), ServiceDiscoveryServiceName())
}

// DeleteEndpointHandler handles DELETE /service/<service_id> request
func DeleteEndpointHandler(c *routing.Context) error {
	// create service discovery object with authenticated user
	serviceDiscovery := ServiceDiscovery(c)

	// Get uri parameter
	serviceID := c.Param("service_id")

	err := serviceDiscovery.DeleteService(
		service.Service{
			Id: bson.ObjectIdHex(serviceID),
		})

	if err != nil {
		return http.NotFound(c, err.Error(), ServiceDiscoveryServiceName())
	}
	return http.Deleted(c, `Service deleted successfuly`, ServiceDiscoveryServiceName())
}

// ManageServiceHandler call management service
func ManageServiceHandler(c *routing.Context) error {
	// create service discovery object with authenticated user
	serviceDiscovery := ServiceDiscovery(c)

	// Get query parameters
	matchingURI := c.QueryArgs().Peek("service")
	managementType := string(c.QueryArgs().Peek("action"))

	service, err := serviceDiscovery.GetEndpointForUri(string(matchingURI))

	if err == nil {
		success, callResponse := service.ServiceManagementCall(managementType)

		if success {
			return http.Ok(c, `{"error": false, "msg": "Service `+managementType+` successfuly.", "service_response": `+strconv.Quote(callResponse)+`}`, ServiceDiscoveryServiceName())
		}
		return http.Error(c, `Service could not be `+managementType+`.`, 400, ServiceDiscoveryServiceName())
	}
	return ServiceNotFound(c)
}

// ManageServiceTypesHandler call management service
func ManageServiceTypesHandler(c *routing.Context) error {
	managementTypesJSON, err := json.Marshal(config.GApiConfiguration.ManagementTypes)

	if err != nil {
		return ServiceNotFound(c)
	}

	return http.Ok(c, string(managementTypesJSON), ServiceDiscoveryServiceName())
}
