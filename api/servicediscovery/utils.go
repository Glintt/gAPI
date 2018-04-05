package servicediscovery

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"encoding/json"
	"errors"
	"strings"
)

func (serviceDisc *ServiceDiscovery) GetEndpointForUri(uri string) (Service, error) {

	if serviceDisc.isService == false {
		resp := http.MakeRequest(config.GET, config.SERVICE_DISCOVERY_URL+config.SERVICE_DISCOVERY_GROUP+"/endpoint?uri="+uri, "", nil)

		if resp.StatusCode() != 200 {
			return Service{}, errors.New("Not found.")
		}

		responseBody := resp.Body()
		var service Service
		json.Unmarshal(responseBody, &service)

		return service, nil

	} else {
		return serviceDisc.FindServiceWithMatchingPrefix(uri)
	}

	return Service{}, errors.New("Not found.")
}

func GetMatchURI(uri string) string {
	uriParts := strings.Split(uri, "/")
	toMatchUri := "/"

	if len(uriParts) > 1 {
		toMatchUri = toMatchUri + uriParts[1]
	} else {
		toMatchUri = uri
	}
	return toMatchUri
}

func (serviceDisc *ServiceDiscovery) FindServiceWithMatchingPrefix(uri string) (Service, error) {
	toMatchUri := GetMatchURI(uri)
	return funcMap[SD_TYPE]["get"].(func(string) (Service, error))(toMatchUri)
}
