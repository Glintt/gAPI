package servicediscovery

import (
	"encoding/json"
	"errors"
	"gAPIManagement/api/servicediscovery/service"
	"gAPIManagement/api/servicediscovery/servicegroup"

	routing "github.com/qiangxue/fasthttp-routing"
)

func ValidateServiceGroupBody(c *routing.Context) (servicegroup.ServiceGroup, error) {
	var s servicegroup.ServiceGroup
	err := json.Unmarshal(c.Request.Body(), &s)

	if s.Name == "" {
		return servicegroup.ServiceGroup{}, errors.New(`{"error" : true, "msg": "Missing body parameters."}`)
	}
	if err != nil {
		return servicegroup.ServiceGroup{}, errors.New(`{"error" : true, "msg": "Error parsing body."}`)
	}

	return s, nil
}

func ValidateServiceBody(c *routing.Context) (service.Service, error) {
	var s service.Service
	err := json.Unmarshal(c.Request.Body(), &s)

	if s.Name == "" || len(s.Hosts) == 0 || s.MatchingURI == "" || s.ToURI == "" || s.APIDocumentation == "" {
		return service.Service{}, errors.New(`{"error" : true, "msg": "Missing body parameters."}`)
	}
	if err != nil {
		return service.Service{}, errors.New(`{"error" : true, "msg": "Error parsing body."}`)
	}

	s.NormalizeService()

	return s, nil
}

func ValidateServiceExists(s service.Service) (service.Service, error) {
	ser, err := sd.FindService(s)

	if err != nil {
		return service.Service{}, errors.New(`{"error":true, "msg":"Resource not found"}`)
	}

	return ser, nil
}
