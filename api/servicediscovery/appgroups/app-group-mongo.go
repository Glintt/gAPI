package appgroups

import (
	"gAPIManagement/api/database"
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/servicediscovery/service"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const SERVICE_APPS_GROUP_COLLECTION = "services_apps_groups"

func CreateApplicationGroupMongo(bodyMap ApplicationGroup) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)
	collection := db.C(SERVICE_APPS_GROUP_COLLECTION)
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)

	bodyMap.Id = bson.NewObjectId()
	err = collection.Insert(&bodyMap)

	database.MongoDBPool.Close(session)

	return err
}

func GetApplicationGroupsMongo(page int, nameFilter string) []ApplicationGroup {
	skips := servicediscovery.PAGE_LENGTH * (page - 1)

	// Get list of application groups
	var groups []ApplicationGroup

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	db.C(SERVICE_APPS_GROUP_COLLECTION).Find(bson.M{"name": bson.RegEx{nameFilter + ".*", ""}}).Sort("name").Skip(skips).Limit(servicediscovery.PAGE_LENGTH).All(&groups)

	database.MongoDBPool.Close(session)

	return groups
}

func GetApplicationGroupByIdMongo(appGroupId string) (ApplicationGroup, error) {
	appGroupIdHex := bson.ObjectIdHex(appGroupId)
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var group ApplicationGroup
	err := db.C(SERVICE_APPS_GROUP_COLLECTION).FindId(appGroupIdHex).One(&group)

	database.MongoDBPool.Close(session)

	return group, err
}

func GetServicesForApplicationGroupMongo(appGroup ApplicationGroup) ([]service.Service, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var servicesList []service.Service
	err := db.C(service.SERVICES_COLLECTION).Find(bson.M{"_id": bson.M{"$in": appGroup.Services}}).All(&servicesList)

	database.MongoDBPool.Close(session)

	return servicesList, err
}

func DeleteApplicationGroupMongo(appGroupId string) error {
	appGroupIdHex := bson.ObjectIdHex(appGroupId)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(SERVICE_APPS_GROUP_COLLECTION).RemoveId(appGroupIdHex)

	database.MongoDBPool.Close(session)

	return err
}

func UpdateApplicationGroupMongo(appGroupId string, newGroup ApplicationGroup) error {
	appGroupIdHex := bson.ObjectIdHex(appGroupId)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	updateGroupQuery := bson.M{"$set": bson.M{"name": newGroup.Name}}
	err := db.C(SERVICE_APPS_GROUP_COLLECTION).UpdateId(appGroupIdHex, updateGroupQuery)

	database.MongoDBPool.Close(session)

	return err
}

func FindServiceApplicationGroupMongo(serviceId string) ApplicationGroup {
	serviceIdHx := bson.ObjectIdHex(serviceId)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var appGroup ApplicationGroup

	query := bson.M{"services": serviceIdHx}
	db.C(SERVICE_APPS_GROUP_COLLECTION).Find(query).One(&appGroup)

	database.MongoDBPool.Close(session)

	return appGroup
}

func AddServiceToGroupMongo(appGroupId string, serviceId string) error {
	serviceGroupIdHex := bson.ObjectIdHex(appGroupId)
	serviceIdHx := bson.ObjectIdHex(serviceId)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHx}}
	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceIdHx}}

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	_, err := db.C(SERVICE_APPS_GROUP_COLLECTION).UpdateAll(bson.M{}, removeFromAllGroups)
	err = db.C(SERVICE_APPS_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, updateGroup)

	database.MongoDBPool.Close(session)
	return err
}

func RemoveServiceFromGroupMongo(appGroupId string, serviceId string) error {
	serviceGroupIdHex := bson.ObjectIdHex(appGroupId)
	serviceIdHx := bson.ObjectIdHex(serviceId)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHx}}

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(SERVICE_APPS_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, removeFromAllGroups)
	database.MongoDBPool.Close(session)

	return err
}

func UngroupedServicesMongo() []service.Service {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var servicesThatMatch []service.Service

	query := []bson.M{
		{"$lookup": bson.M{"from": SERVICE_APPS_GROUP_COLLECTION, "localField": "_id", "foreignField": "services", "as": "service_app_group"}},
		{"$addFields": bson.M{"zeroAppGroups": bson.M{"$not": bson.M{"$size": "$service_app_group"}}}},
		{"$match": bson.M{"zeroAppGroups": true}},
	}

	db.C(service.SERVICES_COLLECTION).Pipe(query).All(&servicesThatMatch)

	database.MongoDBPool.Close(session)

	return servicesThatMatch
}
