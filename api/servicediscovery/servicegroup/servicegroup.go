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
		"list":              GetServiceGroupsMongo,
		"addservicetogroup": AddServiceToGroupMongo,
		"create":            CreateServiceGroupMongo,
		"update":            UpdateServiceGroupMongo,
		"delete":            DeleteServiceGroupMongo,
	},

	"oracle": {
		"list":              GetServiceGroupsOracle,
		"addservicetogroup": AddServiceToGroupOracle,
		"create":            CreateServiceGroupOracle,
		"update":            UpdateServiceGroupOracle,
		"delete":            DeleteServiceGroupOracle,
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
