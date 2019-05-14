package servicediscovery

import (
	"encoding/json"
	"errors"
	"gAPIManagement/api/config"
	"gAPIManagement/api/database"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery/constants"
	"gAPIManagement/api/servicediscovery/service"
	"gAPIManagement/api/servicediscovery/servicegroup"
	sdUtils "gAPIManagement/api/servicediscovery/utils"
	"gAPIManagement/api/utils"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
)

type ServiceDiscovery struct {
	isService          bool
	registeredServices []service.Service
}

var sd ServiceDiscovery

var Methods = map[string]map[string]interface{}{
	"mongo": {
		"delete":        service.DeleteServiceMongo,
		"update":        service.UpdateMongo,
		"create":        service.CreateServiceMongo,
		"list":          service.ListServicesMongo,
		"get":           service.FindMongo,
		"normalize":     service.NormalizeServicesMongo,
		"distincthosts": service.ListAllAvailableHostsMongo},
	"oracle": {
		"delete":        service.DeleteServiceOracle,
		"update":        service.UpdateOracle,
		"create":        service.CreateServiceOracle,
		"list":          service.ListServicesOracle,
		"get":           service.FindOracle,
		"normalize":     service.NormalizeServicesOracle,
		"distincthosts": service.ListAllAvailableHostsOracle},
	"file": {
		"delete":    DeleteServiceFile,
		"update":    UpdateFile,
		"create":    CreateServiceFile,
		"list":      ListServicesFile,
		"get":       FindFile,
		"normalize": NormalizeServicesFile}}

func (serviceDisc *ServiceDiscovery) SetRegisteredServices(rs []service.Service) {
	serviceDisc.registeredServices = rs
}

func ServiceGroupMethods() map[string]interface{} {
	return servicegroup.ServiceGroupMethods[constants.SD_TYPE]
}

func GetServiceDiscoveryObject() *ServiceDiscovery {
	return &sd
}

func InitServiceDiscovery() {
	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" || config.GApiConfiguration.ServiceDiscovery.Type == "oracle" {
		constants.SD_TYPE = config.GApiConfiguration.ServiceDiscovery.Type

		if !database.IsConnectionDone {
			if err := database.InitDatabaseConnection(); err != nil {
				panic(err.Error())
			}
		}
	} else {
		servicesConfig := LoadServicesConfiguration()
		sd.registeredServices = servicesConfig.Services
	}

	sd.isService = true
}

func (service *ServiceDiscovery) IsExternalRequest(requestContxt *routing.Context) bool {
	hosts, _ := Methods[constants.SD_TYPE]["distincthosts"].(func() ([]string, error))()

	requestHost := requestContxt.RemoteIP().String()

	utils.LogMessage("ListAllAvailableHosts = "+strings.Join(hosts, ","), utils.DebugLogType)
	utils.LogMessage("RequestIp = "+requestHost, utils.DebugLogType)

	for _, v := range hosts {
		hostInfo := strings.Split(v, ":")
		if hostInfo[0] == "localhost" {
			hostInfo[0] = "127.0.0.1"
		}
		if hostInfo[0] == requestHost {
			return false
		}
	}
	return true
}

func (service *ServiceDiscovery) SetIsService(isServ bool) {
	service.isService = isServ
}

func (service *ServiceDiscovery) IsService() bool {
	return service.isService
}

func (sd *ServiceDiscovery) GetListOfServicesGroup() ([]servicegroup.ServiceGroup, error) {
	servicesGroup, err := ServiceGroupMethods()["list"].(func() ([]servicegroup.ServiceGroup, error))()

	return servicesGroup, err
}

// func (sd *ServiceDiscovery) AddServiceToGroup(serviceGroupId string, serviceId string) error {
// 	err := ServiceGroupMethods()["addservicetogroup"].(func(string, string) error)(serviceGroupId, serviceId)
// 	return err
// }

func IsServiceReachableFromExternal(service service.Service, sd ServiceDiscovery) bool {
	if !service.UseGroupAttributes || service.GroupId == "" {
		return service.IsReachable
	}

	sgList, err := sd.GetListOfServicesGroup()
	if err != nil {
		return false
	}

	for _, sg := range sgList {
		if sg.Contains(service.Id) {
			return sg.IsReachable
		}
	}
	return false
}

func (serviceDisc *ServiceDiscovery) GetAllServices() ([]service.Service, error) {
	var services []service.Service

	if serviceDisc.isService == false {
		resp := http.MakeRequest(config.GET, config.SERVICE_DISCOVERY_URL+config.SERVICE_DISCOVERY_GROUP+"/services?page=-1", "", nil)

		if resp.StatusCode() != 200 {
			return []service.Service{}, errors.New("Not found.")
		}

		responseBody := resp.Body()
		json.Unmarshal(responseBody, services)
	} else {
		services = Methods[constants.SD_TYPE]["list"].(func(int, string, bool) []service.Service)(-1, "", true)
	}

	return services, nil
}

func (serviceDisc *ServiceDiscovery) GetEndpointForUri(uri string) (service.Service, error) {

	if serviceDisc.isService == false {
		resp := http.MakeRequest(config.GET, config.SERVICE_DISCOVERY_URL+config.SERVICE_DISCOVERY_GROUP+"/endpoint?uri="+uri, "", nil)

		if resp.StatusCode() != 200 {
			return service.Service{}, errors.New("Not found.")
		}

		responseBody := resp.Body()
		var service service.Service
		json.Unmarshal(responseBody, &service)

		return service, nil

	} else {
		service := service.Service{MatchingURI: uri}
		return serviceDisc.FindService(service)
	}
}

func (serviceDisc *ServiceDiscovery) UpdateService(s service.Service) (service.Service, error) {
	_, status := Methods[constants.SD_TYPE]["update"].(func(service.Service, service.Service) (string, int))(s, s)
	if status == 201 {
		return s, nil
	}

	return service.Service{}, errors.New("Not found.")
}

func (serviceDisc *ServiceDiscovery) FindService(s service.Service) (service.Service, error) {
	return Methods[constants.SD_TYPE]["get"].(func(service.Service) (service.Service, error))(s)
}

func (serviceDisc *ServiceDiscovery) FindServiceWithMatchingPrefix(uri string) (service.Service, error) {
	toMatchUri := sdUtils.GetMatchURI(uri)
	return serviceDisc.FindService(service.Service{MatchingURI: toMatchUri})
}
