package servicediscovery

import (
	"encoding/json"
	"errors"
	"gAPIManagement/api/config"
	"io/ioutil"
	"regexp"
)

func UpdateFile(service Service, serviceExists Service) (string, int) {
	var newServices []Service

	for _, element := range sd.registeredServices {
		if element.Id == serviceExists.Id || (element.Name == serviceExists.Name && element.MatchingURI == serviceExists.MatchingURI && element.ToURI == serviceExists.ToURI && element.Domain == serviceExists.Domain) {

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
	s.Id = s.GenerateId()
	sd.registeredServices = append(sd.registeredServices, s)

	go sd.SaveServicesToFile()

	return `{"error" : false, "msg": "Registered service successfuly."}`, 201
}

func ListServicesFile() []Service {
	return sd.registeredServices
}

func DeleteServiceFile(service Service) (string, int) {
	//service, err := FindFile(GetMatchURI(matchingURI))
	service, err := FindFile(service)

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

func FindFile(service Service) (Service, error) {
	for _, rs := range sd.registeredServices {
		//if toMatchUri == rs.MatchingURI {
		//fmt.Println("1=>" + toMatchUri)
		//fmt.Println("2=>" + rs.MatchingURI)
		/*if strings.HasPrefix(toMatchUri, rs.MatchingURI) {
			return rs, nil
		}*/
		//fmt.Println("3=>" + GetMatchingURIRegex(rs.MatchingURI))
		if (rs.MatchingURIRegex == "") {
			rs.MatchingURIRegex = GetMatchingURIRegex(rs.MatchingURI)
		}
		re := regexp.MustCompile(rs.MatchingURIRegex)
		if re.MatchString(service.MatchingURI) {
			return rs, nil
		}
		if rs.Id == service.Id {
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

func NormalizeServicesFile() error{
	var normalizedServices []Service

	for _, rs := range sd.registeredServices {
		rs.NormalizeService()

		normalizedServices = append(normalizedServices, rs)
	}

	sd.registeredServices = normalizedServices
	sd.SaveServicesToFile()
	 
	return nil
}
