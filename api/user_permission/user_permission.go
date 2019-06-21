package user_permission

import (
	"github.com/Glintt/gAPI/api/user_permission/providers"
	"github.com/Glintt/gAPI/api/user_permission/models"
)

const (
	SERVICE_NAME = "USER_PERMISSION"
)

func getRepository() providers.PermissionsRepositoryInterface {
	return providers.GetPermissionRepository()
}

func GetUserPermissions(user_id string) ([]models.UserPermission, error){
	repository := getRepository()
	repository.CreateTransaction()
	
	permissions,err := repository.Get(user_id)
	if err != nil {
		repository.RollbackTransaction()
		return permissions,err
	}
	repository.CommitTransaction()
	return permissions, nil
}

// UserHasPermissionToAccessService check if user has permission with to access service with id = serviceID
func UserHasPermissionToAccessService(userID string, serviceID string) (bool, error){
	repository := getRepository()
	repository.CreateTransaction()
	
	hasPermission, err := repository.HasPermission(userID, serviceID)
	
	repository.CommitTransaction()
	return hasPermission, err
}

func AddPermission(permission models.UserPermission) error {
	repository := getRepository()
	repository.CreateTransaction()

	err := repository.Add(permission)
	if err != nil {
		repository.RollbackTransaction()
		return err
	}
	repository.CommitTransaction()
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
	return nil
}