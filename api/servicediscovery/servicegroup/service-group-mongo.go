package servicegroup

import (
	"gAPIManagement/api/database"

	"gopkg.in/mgo.v2/bson"
)

func GetServiceGroupsMongo() ([]ServiceGroup, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var servicesGroup []ServiceGroup
	err := db.C(SERVICE_GROUP_COLLECTION).Find(bson.M{}).All(&servicesGroup)

	database.MongoDBPool.Close(session)

	return servicesGroup, err
}

func CreateServiceGroupMongo(serviceGroupId string, serviceId string) error {

	return nil
}
