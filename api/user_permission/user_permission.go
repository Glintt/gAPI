package user_permission

import (
	"github.com/Glintt/gAPI/api/user_permission/providers"
	"github.com/Glintt/gAPI/api/user_permission/models"
	"github.com/Glintt/gAPI/api/cache/cacheable"
)

var cacheStorage = cacheable.GetCacheableStorageInstance()

const (
	SERVICE_NAME = "USER_PERMISSION"
)

func getRepository() providers.PermissionsRepositoryInterface {
	return providers.GetPermissionRepository()
}

// GetUserPermissions get user's permissions
func GetUserPermissions(userID string) ([]models.UserPermission, error){
	var permissions []models.UserPermission

	err := cacheStorage.Cacheable("perm_"+userID, func() (interface{}, error) {
		repository := getRepository()
		repository.CreateTransaction()
		
		permissions,err := repository.Get(userID)
		if err != nil {
			repository.RollbackTransaction()
			return permissions, err
		}
		repository.CommitTransaction()

		return permissions, err
	}, &permissions)

	return permissions, err
}

// UserHasPermissionToAccessService check if user has permission with to access service with id = serviceID
func UserHasPermissionToAccessService(userID string, serviceID string) (bool, error){
	var hasPermission bool
	err := cacheStorage.Cacheable("perm_"+userID+"_"+serviceID, func() (interface{}, error) {
		repository := getRepository()
		repository.CreateTransaction()

		hasPerm, err := repository.HasPermission(userID, serviceID)
		
		repository.CommitTransaction()
		return hasPerm, err
	}, &hasPermission)

	return hasPermission, err
}

// AddPermission add permission to a service
func AddPermission(permission models.UserPermission) error {
	repository := getRepository()
	repository.CreateTransaction()

	err := repository.Add(permission)
	if err != nil {
		repository.RollbackTransaction()
		return err
	}
	repository.CommitTransaction()

	cacheStorage.Delete("perm_"+permission.UserId+"_"+permission.ServiceId)
	return nil
}

func UpdatePermission(userId string, permissions []models.UserPermission) error{
	repository := getRepository()	
	repository.CreateTransaction()

	err := repository.Update(userId, permissions)
	if err != nil {
		repository.RollbackTransaction()
		return err
	}
	err = repository.CommitTransaction()

	// Remove cache in order to update instantly results
	cacheStorage.Reset()
	return nil
}

func DeletePermission(userId string, permissionId string) error{
	repository := getRepository()	
	repository.CreateTransaction()

	err := repository.Delete(userId, permissionId)
	if err != nil {
		repository.RollbackTransaction()
		return err
	}
	repository.CommitTransaction()

	// Remove cache in order to update instantly results
	cacheStorage.Delete("perm_"+userId+"_"+permissionId)
	return nil
}
