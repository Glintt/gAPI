package servicegroup

import (
	"errors"
	"gAPIManagement/api/database"
	"gAPIManagement/api/servicediscovery/constants"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetServiceGroupsMongo() ([]ServiceGroup, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var servicesGroup []ServiceGroup
	err := db.C(SERVICE_GROUP_COLLECTION).Find(bson.M{}).All(&servicesGroup)

	database.MongoDBPool.Close(session)

	return servicesGroup, err
}

func CreateServiceGroupMongo(serviceGroup ServiceGroup) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	serviceGroup.Id = bson.NewObjectId()
	collection := db.C(constants.SERVICE_GROUP_COLLECTION)
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		return err
	}
	err = collection.Insert(&serviceGroup)

	database.MongoDBPool.Close(session)

	return err
}

func UpdateServiceGroupMongo(serviceGroupId string, serviceGroup ServiceGroup) error {
	serviceGroupIdObj := bson.ObjectIdHex(serviceGroupId)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(constants.SERVICE_GROUP_COLLECTION).UpdateId(serviceGroupIdObj, serviceGroup)

	database.MongoDBPool.Close(session)
	return err
}

func DeleteServiceGroupMongo(serviceGroupId string) error {
	serviceGroupIdObj := bson.ObjectIdHex(serviceGroupId)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(constants.SERVICE_GROUP_COLLECTION).RemoveId(serviceGroupIdObj)

	_, err = db.C(constants.SERVICES_COLLECTION).UpdateAll(bson.M{}, bson.M{"$set": bson.M{"groupid": nil}})

	database.MongoDBPool.Close(session)

	return err
}

func AddServiceToGroupMongo(serviceGroupId string, serviceId string) error {
	serviceGroupIdHex := bson.ObjectIdHex(serviceGroupId)
	serviceIdHex := bson.ObjectIdHex(serviceId)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHex}}
	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceIdHex}}
	updateService := bson.M{"$set": bson.M{"groupid": serviceGroupIdHex}}

	session, db := database.GetSessionAndDB(database.MONGO_DB)
	err := db.C(constants.SERVICES_COLLECTION).UpdateId(serviceIdHex, updateService)
	if err != nil {
		database.MongoDBPool.Close(session)
		return errors.New("Update Service failed")
	}

	_, err = db.C(constants.SERVICE_GROUP_COLLECTION).UpdateAll(bson.M{}, removeFromAllGroups)
	err = db.C(constants.SERVICE_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, updateGroup)

	database.MongoDBPool.Close(session)

	return err
}
