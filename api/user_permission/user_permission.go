package user_permission

import (
	"github.com/Glintt/gAPI/api/user_permission/providers"
	"github.com/Glintt/gAPI/api/user_permission/models"
)

const (
	SERVICE_NAME = "USER_PERMISSION"
)

func getRepositoryAndBeginTransaction() providers.PermissionsRepositoryInterface {
	repository := providers.GetPermissionRepository()
	repository.CreateTransaction()
	return repository
}

func GetUserPermissions(user_id string) ([]models.UserPermission, error){
	repository := getRepositoryAndBeginTransaction()
	
	permissions,err := repository.Get(user_id)
	if err != nil {
		repository.RollbackTransaction()
		return permissions,err
	}
	repository.CommitTransaction()
	return permissions, nil
}

func AddPermission(permission models.UserPermission) error {
	repository := getRepositoryAndBeginTransaction()

	err := repository.Add(permission)
	if err != nil {
		repository.RollbackTransaction()
		return err
	}
	repository.CommitTransaction()
	return nil
}

func UpdatePermission(userId string, permissions []models.UserPermission) error{
	repository := getRepositoryAndBeginTransaction()	

	err := repository.Update(userId, permissions)
	if err != nil {
		repository.RollbackTransaction()
		return err
	}
	err = repository.CommitTransaction()
	return nil
}

func DeletePermission(userId string, permissionId string) error{
	repository := getRepositoryAndBeginTransaction()	

	err := repository.Delete(userId, permissionId)
	if err != nil {
		repository.RollbackTransaction()
		return err
	}
	repository.CommitTransaction()
	return nil
}
