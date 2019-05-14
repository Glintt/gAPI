package appgroups

import (
	"gopkg.in/mgo.v2/bson"
)

type ApplicationGroup struct {
	Id       bson.ObjectId `bson:"_id" json:"Id"`
	Name     string
	Services []bson.ObjectId
}

var ApplicationGroupMethods = map[string]map[string]interface{}{
	"mongo": {
		"create":                 CreateApplicationGroupMongo,
		"list":                   GetApplicationGroupsMongo,
		"delete":                 DeleteApplicationGroupMongo,
		"getbyid":                GetApplicationGroupByIdMongo,
		"getservicesforappgroup": GetServicesForApplicationGroupMongo,
		"update":                 UpdateApplicationGroupMongo,
		"getappforservice":       FindServiceApplicationGroupMongo,
		"addapptogroup":          AddServiceToGroupMongo,
		"removeappfromgroup":     RemoveServiceFromGroupMongo,
		"ungroupedservices":      UngroupedServicesMongo},

	"oracle": {
		"create":                 CreateApplicationGroupOracle,
		"list":                   GetApplicationGroupsOracle,
		"delete":                 DeleteApplicationGroupOracle,
		"getbyid":                GetApplicationGroupByIdOracle,
		"getservicesforappgroup": GetServicesForApplicationGroupOracle,
		"update":                 UpdateApplicationGroupOracle,
		"getappforservice":       FindServiceApplicationGroupOracle,
		"addapptogroup":          AddServiceToGroupOracle,
		"removeappfromgroup":     RemoveServiceFromGroupOracle,
		"ungroupedservices":      UngroupedServicesOracle},
}
