package servicegroup

import (
	"gopkg.in/mgo.v2/bson"
)

const SERVICE_GROUP_COLLECTION = "services_groups"

type ServiceGroup struct {
	Id           bson.ObjectId `bson:"_id" json:"Id"`
	Name         string
	IsReachable  bool
	HostsEnabled []string
	Services     []bson.ObjectId
}

var ServiceGroupMethods = map[string]map[string]interface{}{
	"mongo": {
		"list":                   GetServiceGroupsMongo,
		"addservicetogroup":      AddServiceToGroupMongo,
		"removeservicefromgroup": RemoveServiceFromGroupMongo,
		"create":                 CreateServiceGroupMongo,
		"update":                 UpdateServiceGroupMongo,
		"delete":                 DeleteServiceGroupMongo,
		"getbyid":                GetServiceGroupByIdMongo,
	},

	"oracle": {
		"list":                   GetServiceGroupsOracle,
		"addservicetogroup":      AddServiceToGroupOracle,
		"removeservicefromgroup": RemoveServiceFromGroupOracle,
		"create":                 CreateServiceGroupOracle,
		"update":                 UpdateServiceGroupOracle,
		"delete":                 DeleteServiceGroupOracle,
		"getbyid":                GetServiceGroupByIdOracle,
	},
}

func (sg *ServiceGroup) Contains(serviceId bson.ObjectId) bool {
	for _, v := range sg.Services {
		if serviceId == v {
			return true
		}
	}
	return false
}
