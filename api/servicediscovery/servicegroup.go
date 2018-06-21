package servicediscovery

import (
	"errors"
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

func (sd *ServiceDiscovery) AddServiceToGroup(serviceGroupId string, serviceId string ) ( error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	serviceGroupIdHex := bson.ObjectIdHex(serviceGroupId)
	serviceIdHex := bson.ObjectIdHex(serviceId)
	
	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHex }}
	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceIdHex }}
	updateService := bson.M{"$set": bson.M{"groupid": serviceGroupIdHex}}

	err := db.C(SERVICES_COLLECTION).UpdateId(serviceIdHex, updateService)
	if err != nil {
		database.MongoDBPool.Close(session)
		return errors.New("Update Service failed")
	}

	_,err = db.C(SERVICE_GROUP_COLLECTION).UpdateAll(bson.M{}, removeFromAllGroups)
	err = db.C(SERVICE_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, updateGroup)

	database.MongoDBPool.Close(session)
	return nil
}