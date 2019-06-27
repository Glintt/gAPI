package servicegroup

import (
	"errors"

	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const SERVICE_GROUP_COLLECTION = "services_groups"

type ServiceGroupMongoRepository struct {
	Session    *mgo.Session
	Db         *mgo.Database
	Collection *mgo.Collection
}

// OpenTransaction open new database transaction
func (sgmr *ServiceGroupMongoRepository) OpenTransaction() error {
	return nil
}

// CommitTransaction commit database transaction
func (sgmr *ServiceGroupMongoRepository) CommitTransaction() {}

// RollbackTransaction rollback database transaction
func (sgmr *ServiceGroupMongoRepository) RollbackTransaction() {}

// Release release database transaction to the pool
func (sgmr *ServiceGroupMongoRepository) Release() {
	database.MongoDBPool.Close(sgmr.Session)
}

// GetServiceGroups get list of service groups
func (sgmr *ServiceGroupMongoRepository) GetServiceGroups() ([]ServiceGroup, error) {
	var servicesGroup []ServiceGroup
	err := sgmr.Collection.Find(bson.M{}).All(&servicesGroup)
	return servicesGroup, err
}

// GetServiceGroupById get service group by id
func (sgmr *ServiceGroupMongoRepository) GetServiceGroupById(serviceGroup string) (ServiceGroup, error) {
	var sg ServiceGroup
	var err error
	if bson.IsObjectIdHex(serviceGroup) {
		err = sgmr.Collection.FindId(bson.ObjectIdHex(serviceGroup)).One(&sg)
	} else {
		err = sgmr.Collection.Find(bson.M{"name": serviceGroup}).One(&sg)
	}

	return sg, err
}

// CreateServiceGroup create new service group
func (sgmr *ServiceGroupMongoRepository) CreateServiceGroup(serviceGroup ServiceGroup) error {
	serviceGroup.Id = bson.NewObjectId()
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := sgmr.Collection.EnsureIndex(index)
	if err != nil {
		return err
	}
	return sgmr.Collection.Insert(&serviceGroup)
}

// UpdateServiceGroup update an existing service group
func (sgmr *ServiceGroupMongoRepository) UpdateServiceGroup(serviceGroupID string, serviceGroup ServiceGroup) error {
	serviceGroupIdObj := bson.ObjectIdHex(serviceGroupID)

	return sgmr.Collection.UpdateId(serviceGroupIdObj, serviceGroup)
}

// DeleteServiceGroup delete an existing service group
func (sgmr *ServiceGroupMongoRepository) DeleteServiceGroup(serviceGroupID string) error {
	serviceGroupIdObj := bson.ObjectIdHex(serviceGroupID)

	err := sgmr.Collection.RemoveId(serviceGroupIdObj)

	_, err = sgmr.Db.C(constants.SERVICES_COLLECTION).UpdateAll(bson.M{}, bson.M{"$set": bson.M{"groupid": nil}})
	return err
}

// AddServiceToGroup add service to an existing service group
func (sgmr *ServiceGroupMongoRepository) AddServiceToGroup(serviceGroupID string, serviceID string) error {
	serviceGroupIdHex := bson.ObjectIdHex(serviceGroupID)
	serviceIdHex := bson.ObjectIdHex(serviceID)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHex}}
	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceIdHex}}
	updateService := bson.M{"$set": bson.M{"groupid": serviceGroupIdHex}}

	err := sgmr.Db.C(constants.SERVICES_COLLECTION).UpdateId(serviceIdHex, updateService)
	if err != nil {
		return errors.New("Update Service failed")
	}

	_, err = sgmr.Collection.UpdateAll(bson.M{}, removeFromAllGroups)
	if err != nil {
		return errors.New("Update Service failed")
	}
	return sgmr.Collection.UpdateId(serviceGroupIdHex, updateGroup)
}

// RemoveServiceFromGroup remove service from an existing service group
func (sgmr *ServiceGroupMongoRepository) RemoveServiceFromGroup(serviceGroupID string, serviceID string) error {
	serviceGroupIdHex := bson.ObjectIdHex(serviceGroupID)
	serviceIdHex := bson.ObjectIdHex(serviceID)

	updateGroup := bson.M{"$pull": bson.M{"services": serviceIdHex}}
	updateService := bson.M{"$set": bson.M{"groupid": nil, "usegroupattributes": false}}

	err := sgmr.Db.C(constants.SERVICES_COLLECTION).UpdateId(serviceIdHex, updateService)

	if err != nil {
		return err
	}

	return sgmr.Collection.UpdateId(serviceGroupIdHex, updateGroup)
}
