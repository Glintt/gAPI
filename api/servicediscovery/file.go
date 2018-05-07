package servicediscovery

import (
	"gAPIManagement/api/config"
	"encoding/json"
	"errors"
	"io/ioutil"
)

func UpdateFile(service Service, serviceExists Service) (string, int) {
	var newServices []Service

	for _, element := range sd.registeredServices {
		if element.Name == serviceExists.Name && element.MatchingURI == serviceExists.MatchingURI && element.ToURI == serviceExists.ToURI && element.Domain == serviceExists.Domain {

		} else {
			newServices = append(newServices, element)
		}
	}

	newServices = append(newServices, service)
	sd.registeredServices = newServices

	go sd.SaveServicesToFile()

	return `{"error" : false, "msg": "Service updated successfuly."}`, 201
}

func CreateServiceFile(s Service) (string, int) {
	sd.registeredServices = append(sd.registeredServices, s)

	go sd.SaveServicesToFile()

	return `{"error" : false, "msg": "Registered service successfuly."}`, 201
}

func ListServicesFile(page int) []Service {
	if page == -1 {
		return sd.registeredServices
	}
	
	from, to := pageFromTo(page, len(sd.registeredServices))

	return sd.registeredServices[from:to]
}

func DeleteServiceFile(matchingURI string) (string, int) {
	service, err := FindFile(GetMatchURI(matchingURI))

	if err != nil {
		return `{"error": true, "msg": "Not found"}`, 404
	}

	var newServices []Service

	for _, element := range sd.registeredServices {
		if element.Name == service.Name && element.MatchingURI == service.MatchingURI && element.ToURI == service.ToURI && element.Domain == service.Domain {
		} else {
			newServices = append(newServices, element)
		}
	}

	sd.registeredServices = newServices

	go sd.SaveServicesToFile()

	return `{"error": false, "msg": "Removed successfully."}`, 200
}

func FindFile(toMatchUri string) (Service, error) {
	for _, rs := range sd.registeredServices {
		if toMatchUri == rs.MatchingURI {
			return rs, nil
		}
	}
	return Service{}, errors.New("Not found.")
}

func (service *ServiceDiscovery) SaveServicesToFile() {
	var reg map[string][]Service
	reg = make(map[string][]Service)
	reg["services"] = service.registeredServices

	registeredServicesJson, err := json.Marshal(reg)

	if err != nil {
		return
	}

	err = ioutil.WriteFile(config.CONFIGS_LOCATION+config.SERVICE_DISCOVERY_CONFIG_FILE, registeredServicesJson, 0777)
}
