package providers

import (
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/user_permission/models"
)

type PermissionsRepositoryInterface interface{	
	Get(userId string) ([]models.UserPermission, error)
	Add(userPermission models.UserPermission) error
	Update(userId string, userPermission []models.UserPermission) error
	Delete(userId string, serviceId string) error
	DeleteAll(userId string) error 
	CreateTransaction() error
	CommitTransaction() error
	RollbackTransaction() error
}

func GetPermissionRepository() PermissionsRepositoryInterface {
	// if constants.SD_TYPE == "mongo" {
	// 	return &ServiceMongoRepository{
	// 		User: user,
	// 	}
	// }
	if database.SD_TYPE == "oracle" {
		repos := &PermissionsOracleRepository{}
		repos.Init()
		return repos
	}
	return nil
}