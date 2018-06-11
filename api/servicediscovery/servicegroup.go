package servicediscovery

import (
	"gAPIManagement/api/database"
	"gopkg.in/mgo.v2/bson"
)


type ServiceGroup struct {
	Id                    bson.ObjectId `bson:"_id" json:"Id"`
	Name string
	IsReachable bool
	HostsEnabled []string
	Services []bson.ObjectId
}

func (sg *ServiceGroup) Contains(s Service) bool {
	for _, v := range sg.Services {
		if s.Id == v {
			return true
		}
	}
	return false
}

func (sd *ServiceDiscovery) GetListOfServicesGroup() ([]ServiceGroup, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var servicesGroup []ServiceGroup
	err := db.C(SERVICE_GROUP_COLLECTION).Find(bson.M{}).All(&servicesGroup)

	database.MongoDBPool.Close(session)
	
	return servicesGroup, err
}