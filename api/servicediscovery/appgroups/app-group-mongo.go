package appgroups

import (
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/servicediscovery/service"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	userModels "github.com/Glintt/gAPI/api/users/models"
)

type AppGroupMongoRepository struct {
	Session    *mgo.Session
	Db         *mgo.Database
	Collection *mgo.Collection
	User userModels.User
}

// OpenTransaction open new database transaction
func (agmr *AppGroupMongoRepository) OpenTransaction() error { return nil }

// CommitTransaction commit database transaction
func (agmr *AppGroupMongoRepository) CommitTransaction() {}

// RollbackTransaction rollback current transaction
func (agmr *AppGroupMongoRepository) RollbackTransaction() {}

// Release release current database connection
func (agmr *AppGroupMongoRepository) Release() {
	database.MongoDBPool.Close(agmr.Session)
}

// CreateApplicationGroup create application group
func (agmr *AppGroupMongoRepository) CreateApplicationGroup(bodyMap ApplicationGroup) error {
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := agmr.Collection.EnsureIndex(index)
	if err != nil {
		return err
	}
	bodyMap.Id = bson.NewObjectId()
	return agmr.Collection.Insert(&bodyMap)
}

// GetApplicationGroupsForUser Gets application groups an user has access to
func (agmr *AppGroupMongoRepository) GetApplicationGroupsForUser(userID string) ([]ApplicationGroup, error) {
	return agmr.GetApplicationGroups(-1, ""), nil
}


// GetApplicationGroups get list of application groups
func (agmr *AppGroupMongoRepository) GetApplicationGroups(page int, nameFilter string) []ApplicationGroup {
	skips := constants.PAGE_LENGTH * (page - 1)

	var groups []ApplicationGroup

	agmr.Collection.Find(bson.M{"name": bson.RegEx{nameFilter + ".*", ""}}).Sort("name").Skip(skips).Limit(constants.PAGE_LENGTH).All(&groups)

	return groups
}

// GetApplicationGroupByID get application group by id
func (agmr *AppGroupMongoRepository) GetApplicationGroupByID(appGroupID string) (ApplicationGroup, error) {
	appGroupIDHex := bson.ObjectIdHex(appGroupID)

	var group ApplicationGroup
	err := agmr.Collection.FindId(appGroupIDHex).One(&group)

	return group, err
}

// GetServicesForApplicationGroup get application group's services
func (agmr *AppGroupMongoRepository) GetServicesForApplicationGroup(appGroup ApplicationGroup) ([]service.Service, error) {
	var servicesList []service.Service
	err := agmr.Collection.Find(bson.M{"_id": bson.M{"$in": appGroup.Services}}).All(&servicesList)

	return servicesList, err
}

// DeleteApplicationGroup delete application group by id
func (agmr *AppGroupMongoRepository) DeleteApplicationGroup(appGroupID string) error {
	appGroupIDHex := bson.ObjectIdHex(appGroupID)

	return agmr.Collection.RemoveId(appGroupIDHex)
}

// UpdateApplicationGroup update application group with id appGroupID
func (agmr *AppGroupMongoRepository) UpdateApplicationGroup(appGroupID string, newGroup ApplicationGroup) error {
	appGroupIDHex := bson.ObjectIdHex(appGroupID)

	updateGroupQuery := bson.M{"$set": bson.M{"name": newGroup.Name}}
	return agmr.Collection.UpdateId(appGroupIDHex, updateGroupQuery)
}

// FindServiceApplicationGroup get application group for service with id serviceID
func (agmr *AppGroupMongoRepository) FindServiceApplicationGroup(serviceID string) ApplicationGroup {
	serviceIDHx := bson.ObjectIdHex(serviceID)

	var appGroup ApplicationGroup

	query := bson.M{"services": serviceIDHx}
	agmr.Collection.Find(query).One(&appGroup)

	return appGroup
}

// AddServiceToGroup add service with id serviceID to application group with id appGroupID
func (agmr *AppGroupMongoRepository) AddServiceToGroup(appGroupID string, serviceID string) error {
	serviceGroupIDHex := bson.ObjectIdHex(appGroupID)
	serviceIDHx := bson.ObjectIdHex(serviceID)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIDHx}}
	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceIDHx}}

	_, err := agmr.Collection.UpdateAll(bson.M{}, removeFromAllGroups)
	if err != nil {
		return err
	}
	return agmr.Collection.UpdateId(serviceGroupIDHex, updateGroup)
}

// RemoveServiceFromGroup remove service from application group
func (agmr *AppGroupMongoRepository) RemoveServiceFromGroup(appGroupID string, serviceID string) error {
	serviceGroupIDHex := bson.ObjectIdHex(appGroupID)
	serviceIDHx := bson.ObjectIdHex(serviceID)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIDHx}}

	return agmr.Collection.UpdateId(serviceGroupIDHex, removeFromAllGroups)
}

// UngroupedServices get list of services without any application group
func (agmr *AppGroupMongoRepository) UngroupedServices() []service.Service {
	var servicesThatMatch []service.Service

	query := []bson.M{
		{"$lookup": bson.M{"from": constants.SERVICE_APPS_GROUP_COLLECTION, "localField": "_id", "foreignField": "services", "as": "service_app_group"}},
		{"$addFields": bson.M{"zeroAppGroups": bson.M{"$not": bson.M{"$size": "$service_app_group"}}}},
		{"$match": bson.M{"zeroAppGroups": true}},
	}

	agmr.Collection.Pipe(query).All(&servicesThatMatch)

	return servicesThatMatch
}
