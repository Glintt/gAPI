package servicediscovery

import (
	"strings"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	"github.com/Glintt/gAPI/api/servicediscovery/servicegroup"
	sdUtils "github.com/Glintt/gAPI/api/servicediscovery/utils"
	"github.com/Glintt/gAPI/api/users"
	"github.com/Glintt/gAPI/api/utils"
	routing "github.com/qiangxue/fasthttp-routing"
)

// ServiceDiscovery object with service discovery logic
type ServiceDiscovery struct {
	registeredServices []service.Service
	User               users.User
}

// ServiceGroupMethods methods to access service group's repository
func ServiceGroupMethods() map[string]interface{} {
	return servicegroup.ServiceGroupMethods[database.SD_TYPE]
}

// GetServiceDiscoveryObject return service discovery object with request user context 
func GetServiceDiscoveryObject(user users.User) *ServiceDiscovery {
	return &ServiceDiscovery{
		User: user,
	}
}

// GetInternalServiceDiscoveryObject return service discovery object with internal user 
func GetInternalServiceDiscoveryObject() *ServiceDiscovery {
	user := users.GetInternalAPIUser()
	return &ServiceDiscovery{
		User: user,
	}
}

// InitServiceDiscovery initializes all components required for service discovery to work
func InitServiceDiscovery() {
	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" || config.GApiConfiguration.ServiceDiscovery.Type == "oracle" {
		database.SD_TYPE = config.GApiConfiguration.ServiceDiscovery.Type

		if !database.IsConnectionDone {
			if err := database.InitDatabaseConnection(); err != nil {
				panic(err.Error())
			}
		}
	}
}

// IsServiceReachableFromExternal check if service is reachable from external
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

func (serviceDisc *ServiceDiscovery) SetRegisteredServices(rs []service.Service) {
	serviceDisc.registeredServices = rs
}

func (s *ServiceDiscovery) IsExternalRequest(requestContxt *routing.Context) bool {
	hosts, _ := service.GetServicesRepository(s.User).ListAllAvailableHosts()

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

func (sd *ServiceDiscovery) GetListOfServicesGroup() ([]servicegroup.ServiceGroup, error) {
	servicesGroup, err := ServiceGroupMethods()["list"].(func() ([]servicegroup.ServiceGroup, error))()

	return servicesGroup, err
}

func (serviceDisc *ServiceDiscovery) GetAllServices() ([]service.Service, error) {
	var services []service.Service
	services = service.GetServicesRepository(serviceDisc.User).ListServices(-1, "")

	return services, nil
}

func (serviceDisc *ServiceDiscovery) ListServices(page int, searchQuery string) []service.Service {
	return service.GetServicesRepository(serviceDisc.User).ListServices(page, searchQuery)
}

func (serviceDisc *ServiceDiscovery) DeleteService(s service.Service) error {
	return service.GetServicesRepository(serviceDisc.User).DeleteService(s)
}

func (serviceDisc *ServiceDiscovery) GetEndpointForUri(uri string) (service.Service, error) {
	service := service.Service{MatchingURI: uri}
	return serviceDisc.FindService(service)
}

func (serviceDisc *ServiceDiscovery) NormalizeServices() error {
	return service.GetServicesRepository(serviceDisc.User).NormalizeServices()
}

func (serviceDisc *ServiceDiscovery) CreateService(s service.Service) (service.Service, error) {
	return service.GetServicesRepository(serviceDisc.User).CreateService(s)
}

func (serviceDisc *ServiceDiscovery) UpdateService(s service.Service) (service.Service, error) {
	_, err := service.GetServicesRepository(serviceDisc.User).Update(s, s)
	if err == nil {
		return s, nil
	}

	return service.Service{}, err
}

func (serviceDisc *ServiceDiscovery) FindByIdentifierOrMatchingUri(identifier string, matchingURI string) (service.Service, error) {
	if identifier != "" {
		return serviceDisc.FindService(
			service.Service{
				Identifier: string(identifier),
			})
	}else {
		return serviceDisc.GetEndpointForUri(string(matchingURI))
	}
}

func (serviceDisc *ServiceDiscovery) FindService(s service.Service) (service.Service, error) {
	return service.GetServicesRepository(serviceDisc.User).Find(s)
}

func (serviceDisc *ServiceDiscovery) FindServiceWithMatchingPrefix(uri string) (service.Service, error) {
	toMatchUri := sdUtils.GetMatchURI(uri)
	return serviceDisc.FindService(service.Service{MatchingURI: toMatchUri})
}
