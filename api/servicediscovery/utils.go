package servicediscovery

import (
	"encoding/json"
	"errors"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"regexp"
)

func (serviceDisc *ServiceDiscovery) GetAllServices() ([]Service, error) {

	if serviceDisc.isService == false {
		resp := http.MakeRequest(config.GET, config.SERVICE_DISCOVERY_URL+config.SERVICE_DISCOVERY_GROUP+"/services", "", nil)

		if resp.StatusCode() != 200 {
			return []Service{}, errors.New("Not found.")
		}

		responseBody := resp.Body()
		var services []Service
		json.Unmarshal(responseBody, services)

		return services, nil

	} else {
		services := funcMap[SD_TYPE]["list"].(func() []Service)()
		return services, nil
	}

	//return []Service{}, errors.New("Not found.")
}

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
		service := Service{MatchingURI: uri}
		return serviceDisc.FindService(service)
	}

	//return Service{}, errors.New("Not found.")
}

/*func GetMatchURI(uri string) string {
	f := func(c rune) bool {
		return c == '/'
	}
	uriParts := strings.FieldsFunc(uri, f)

	toMatchUri := "/" + strings.Join(uriParts, "/") + "/"

	return toMatchUri
}*/

func GetMatchingURIRegex(uri string) string {
	s := uri
	re := regexp.MustCompile("^(\\^/)?/?")
	s = re.ReplaceAllString(s, "^/")
	re = regexp.MustCompile("(/(\\.\\*)?)?$")
	s = re.ReplaceAllString(s, "((/.*)|$)")
	return s
}

func (serviceDisc *ServiceDiscovery) FindService(service Service) (Service, error) {
	//toMatchUri := uri
	//toMatchUri := GetMatchURI(uri)
	return funcMap[SD_TYPE]["get"].(func(Service) (Service, error))(service)
}