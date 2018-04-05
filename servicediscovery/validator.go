package servicediscovery

import (
	"encoding/json"
	"errors"

	routing "github.com/qiangxue/fasthttp-routing"
)

func ValidateServiceBody(c *routing.Context) (Service, error) {
	var s Service
	err := json.Unmarshal(c.Request.Body(), &s)

	if s.Name == "" || s.Domain == "" || s.Port == "" || s.MatchingURI == "" || s.ToURI == "" || s.APIDocumentation == "" {
		return Service{}, errors.New(`{"error" : true, "msg": "Missing body parameters."}`)
	}
	if err != nil {
		return Service{}, errors.New(`{"error" : true, "msg": "Error parsing body."}`)
	}

	return s, nil
}

func ValidateServiceExists(s Service) (Service, error) {
	service, err := sd.FindServiceWithMatchingPrefix(s.MatchingURI)

	if err != nil {
		return Service{}, errors.New(`{"error":true, "msg":"Resource not found"}`)
	}

	return service, nil
}
