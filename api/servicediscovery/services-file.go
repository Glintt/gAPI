package servicediscovery

import (
	"encoding/json"
	"errors"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	sdUtils "github.com/Glintt/gAPI/api/servicediscovery/utils"
	"github.com/Glintt/gAPI/api/utils"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type ServicesConfig struct {
	Services []service.Service `json:"services"`
}

func LoadServicesConfiguration() ServicesConfig {
	servicesJSON, err := utils.LoadJsonFile(config.CONFIGS_LOCATION + config.SERVICE_DISCOVERY_CONFIG_FILE)

	if err != nil {
		return ServicesConfig{}
	}

	var sc ServicesConfig
	json.Unmarshal(servicesJSON, &sc)
	return sc
}

func UpdateFile(s service.Service, serviceExists service.Service) (string, int) {
	var newServices []service.Service

	for _, element := range sd.registeredServices {
		if element.Id == serviceExists.Id || (element.Name == serviceExists.Name && element.MatchingURI == serviceExists.MatchingURI && element.ToURI == serviceExists.ToURI && element.Domain == serviceExists.Domain) {

		} else {
			newServices = append(newServices, element)
		}
	}

	newServices = append(newServices, s)
	sd.registeredServices = newServices

	go sd.SaveServicesToFile()

	return `{"error" : false, "msg": "Service updated successfuly."}`, 201
}

func CreateServiceFile(s service.Service) (string, int) {
	s.Id = s.GenerateId()
	sd.registeredServices = append(sd.registeredServices, s)

	go sd.SaveServicesToFile()

	return `{"error" : false, "msg": "Registered service successfuly."}`, 201
}

func ListServicesFile(page int, filterQuery string) []service.Service {
	var servicesList []service.Service
	if filterQuery != "" {
		for _, v := range sd.registeredServices {
			if strings.Contains(strings.ToLower(v.Name), strings.ToLower(filterQuery)) || strings.Contains(strings.ToLower(v.MatchingURI), strings.ToLower(filterQuery)) {
				servicesList = append(servicesList, v)
			}
		}
	} else {
		servicesList = sd.registeredServices
	}

	sort.Slice(servicesList, func(i, j int) bool { return servicesList[i].MatchingURI < servicesList[j].MatchingURI })

	if page == -1 {
		return servicesList
	}
	from, to := database.PageFromTo(page, constants.PAGE_LENGTH, len(servicesList))

	return servicesList[from:to]
}

func DeleteServiceFile(s service.Service) (string, int) {
	//service, err := FindFile(GetMatchURI(matchingURI))
	s, err := FindFile(s)

	if err != nil {
		return `{"error": true, "msg": "Not found"}`, 404
	}

	var newServices []service.Service

	for _, element := range sd.registeredServices {
		if element.Name == s.Name && element.MatchingURI == s.MatchingURI && element.ToURI == s.ToURI {

		} else {
			newServices = append(newServices, element)
		}
	}

	sd.registeredServices = newServices

	go sd.SaveServicesToFile()

	return `{"error": false, "msg": "Removed successfully."}`, 200
}

func FindFile(s service.Service) (service.Service, error) {
	for _, rs := range sd.registeredServices {
		if rs.MatchingURIRegex == "" {
			rs.MatchingURIRegex = sdUtils.GetMatchingURIRegex(rs.MatchingURI)
		}
		re := regexp.MustCompile(rs.MatchingURIRegex)
		if re.MatchString(s.MatchingURI) {
			return rs, nil
		}
		if rs.Id == s.Id {
			return rs, nil
		}
	}
	return service.Service{}, errors.New("Not found.")
}

func (s *ServiceDiscovery) SaveServicesToFile() {
	var reg map[string][]service.Service
	reg = make(map[string][]service.Service)
	reg["services"] = s.registeredServices

	registeredServicesJson, err := json.Marshal(reg)

	if err != nil {
		return
	}

	err = ioutil.WriteFile(config.CONFIGS_LOCATION+config.SERVICE_DISCOVERY_CONFIG_FILE, registeredServicesJson, 0777)
}

func NormalizeServicesFile() error {
	var normalizedServices []service.Service

	for _, rs := range sd.registeredServices {
		rs.NormalizeService()

		normalizedServices = append(normalizedServices, rs)
	}

	sd.registeredServices = normalizedServices
	sd.SaveServicesToFile()

	return nil
}
