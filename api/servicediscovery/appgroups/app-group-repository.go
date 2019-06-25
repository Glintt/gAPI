package appgroups

import (
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	userModels "github.com/Glintt/gAPI/api/users/models"
)

type AppGroupRepository interface {
	UpdateApplicationGroup(appGroupID string, newGroup ApplicationGroup) error
	FindServiceApplicationGroup(serviceID string) ApplicationGroup
	CreateApplicationGroup(bodyMap ApplicationGroup) error
	GetApplicationGroups(page int, nameFilter string) []ApplicationGroup
	GetApplicationGroupByID(appGroupID string) (ApplicationGroup, error)
	GetServicesForApplicationGroup(appGroup ApplicationGroup) ([]service.Service, error)
	DeleteApplicationGroup(appGroupID string) error
	AddServiceToGroup(appGroupID string, serviceID string) error
	RemoveServiceFromGroup(appGroupID string, serviceID string) error
	UngroupedServices() []service.Service
	GetApplicationGroupsForUser(userID string) ([]ApplicationGroup, error)

	OpenTransaction() error
	CommitTransaction()
	RollbackTransaction()
	Release()
}

// NewAppGroupRepository create an application group repository based on the database
func NewAppGroupRepository(user userModels.User) AppGroupRepository {
	if database.SD_TYPE == "mongo" {
		session, db := database.GetSessionAndDB(database.MONGO_DB)
		collection := db.C(constants.SERVICE_APPS_GROUP_COLLECTION)

		return &AppGroupMongoRepository{
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
		return &AppGroupOracleRepository{
			Db:      db,
			DbError: err,
		}
	}
	return nil
}
