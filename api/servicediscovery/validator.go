package servicediscovery

import (
	"encoding/json"
	"errors"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	"github.com/Glintt/gAPI/api/servicediscovery/servicegroup"

	"gopkg.in/mgo.v2/bson"
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
// ValidateServiceBody validate service request body
func ValidateServiceBody(c *routing.Context) (service.Service, error) {
	var s service.Service
	err := json.Unmarshal(c.Request.Body(), &s)

	if s.Name == "" || len(s.Hosts) == 0 || s.MatchingURI == "" || s.ToURI == "" || s.APIDocumentation == "" {
		return service.Service{}, errors.New("Missing body parameters")
	}
	if err != nil {
		return service.Service{}, errors.New(`Error parsing body`)
	}

	s.NormalizeService()

	return s, nil
}

// ValidateServiceExists validate if service already exists 
func ValidateServiceExists(s service.Service) (service.Service, error) {
	sd := GetInternalServiceDiscoveryObject()
	ser, err := sd.FindService(s)

	if err != nil {
		return service.Service{}, errors.New(`Service not found`)
	}

	return ser, nil
}


// ValidateServiceBodyAndServiceExists validates both request body and if service already exists 
func ValidateServiceBodyAndServiceExists(c *routing.Context, serviceID string) (service.Service, int, error) {
	s, err := ValidateServiceBody(c)
	if err != nil {
		return service.Service{}, 400, err
	}

	if serviceID != "" {
		s.Id = bson.ObjectIdHex(serviceID)
	}

	_, err = ValidateServiceExists(s)
	if err != nil {
		return service.Service{}, 404, err
	}
	return s, 200, nil
}