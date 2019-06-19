package servicegroup

import (
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/users"
)

type ServiceGroupRepository interface {
	GetServiceGroups() ([]ServiceGroup, error)
	GetServiceGroupById(serviceGroup string) (ServiceGroup, error)
	AddServiceToGroup(serviceGroupId string, serviceId string) error
	RemoveServiceFromGroup(serviceGroupId string, serviceId string) error
	CreateServiceGroup(serviceGroup ServiceGroup) error
	UpdateServiceGroup(serviceGroupId string, serviceGroup ServiceGroup) error
	DeleteServiceGroup(serviceGroupId string) error

	OpenTransaction() error
	CommitTransaction()
	RollbackTransaction()
	Release()
}

// NewServiceGroupRepository create an application group repository based on the database
func NewServiceGroupRepository(user users.User) ServiceGroupRepository {
	if database.SD_TYPE == "mongo" {
		session, db := database.GetSessionAndDB(database.MONGO_DB)
		collection := db.C(SERVICE_GROUP_COLLECTION)

		return &ServiceGroupMongoRepository{
			Session:    session,
			Db:         db,
			Collection: collection,
		}
	}
	if database.SD_TYPE == "oracle" {
		db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
		if err != nil {
			return nil
		}
		return &ServiceGroupOracleRepository{
			Db:      db,
			DbError: err,
		}
	}
	return nil
}