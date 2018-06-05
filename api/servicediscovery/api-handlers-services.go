package servicediscovery

import (
	"gAPIManagement/api/config"
	"encoding/json"
	"strconv"
	"gAPIManagement/api/http"
	"github.com/qiangxue/fasthttp-routing"
)

func ServiceNotFound(c *routing.Context) error {
	http.Response(c, `{"error":true, "msg": "Not found."}`, 404, SERVICE_NAME)
	return nil
}

func NormalizeServices(c *routing.Context) error {
	err := Methods[SD_TYPE]["normalize"].(func() (error))()
	if err != nil {
		http.Response(c, `{"error":true, "msg": "Normalization failed."}`, 400, SERVICE_NAME)
		return err
	}
	http.Response(c, `{"error":false, "msg": "Normalization done."}`, 200, SERVICE_NAME)
	return nil
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

	service.MatchingURIRegex = GetMatchingURIRegex(service.MatchingURI)
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
			http.Response(c, `{"error": false, "msg": "Service ` + managementType + ` successfuly.", "service_response": ` + strconv.Quote(callResponse) + `}` , 200, SERVICE_NAME)
			return nil
		}
		http.Response(c, `{"error": true, "msg": "Service could not be ` + managementType + `."}`, 400, SERVICE_NAME)
		return nil
	}
	http.Response(c, `{"error": true, "msg": "Not found"}`, 404, SERVICE_NAME)
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
	
	http.Response(c, response, statusCode, SERVICE_NAME)
	return nil
}