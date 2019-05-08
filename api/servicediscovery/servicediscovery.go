package servicediscovery

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/database"
	"gAPIManagement/api/utils"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
)

type ServiceDiscovery struct {
	isService          bool
	registeredServices []Service
}

var sd ServiceDiscovery

var SERVICE_NAME = "/service-discovery"
var PAGE_LENGTH = 10
var SD_TYPE = "file"

var Methods = map[string]map[string]interface{}{
	"mongo": {
		"delete":    DeleteServiceMongo,
		"update":    UpdateMongo,
		"create":    CreateServiceMongo,
		"list":      ListServicesMongo,
		"get":       FindMongo,
		"normalize": NormalizeServicesMongo},
	"oracle": {
		"delete":    DeleteServiceOracle,
		"update":    UpdateOracle,
		"create":    CreateServiceOracle,
		"list":      ListServicesOracle,
		"get":       FindOracle,
		"normalize": NormalizeServicesOracle},
	"file": {
		"delete":    DeleteServiceFile,
		"update":    UpdateFile,
		"create":    CreateServiceFile,
		"list":      ListServicesFile,
		"get":       FindFile,
		"normalize": NormalizeServicesFile}}

func (serviceDisc *ServiceDiscovery) SetRegisteredServices(rs []Service) {
	serviceDisc.registeredServices = rs
}

func GetServiceDiscoveryObject() *ServiceDiscovery {
	return &sd
}

func InitServiceDiscovery() {
	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" || config.GApiConfiguration.ServiceDiscovery.Type == "oracle" {
		SD_TYPE = config.GApiConfiguration.ServiceDiscovery.Type

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
	hosts, _ := ListAllAvailableHosts()
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
