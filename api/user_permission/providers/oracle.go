package providers

import (
	"github.com/Glintt/gAPI/api/user_permission/models"
	"github.com/Glintt/gAPI/api/database"
	"errors"
	"database/sql"
	"fmt"
)

type PermissionsOracleRepository struct {
	tx *sql.Tx
	db *sql.DB
}

const (
	GET_PERMISSIONS_QUERY = `
select user_id, service_id from gapi_user_services_permissions where user_id = :user_id
`
	DELETE_PERMISSION_QUERY = `delete from gapi_user_services_permissions where user_id = :user_id and service_id = :service_id` 
	DELETE_ALL_USER_PERMISSIONS_QUERY = `delete from gapi_user_services_permissions where user_id = :user_id` 
	CREATE_USER_PERMISSION_QUERY = `insert into gapi_user_services_permissions(user_id, service_id) values (:user_id, :service_id)` 
)

func (por *PermissionsOracleRepository) CreateTransaction() error {
	tx, err := por.db.Begin()
	if err != nil {
		return errors.New("Error creating transaction " + err.Error())
	}
	por.tx = tx
	return nil
}

func (por *PermissionsOracleRepository) CommitTransaction() error {
	err := por.tx.Commit()
	return err
}

func (por *PermissionsOracleRepository) RollbackTransaction() error {
	por.tx.Rollback()
	return nil
}

func (por *PermissionsOracleRepository) Init() error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return errors.New("Error connecting to database: " + err.Error())
	}	
	por.db = db
	return nil
}

func (por *PermissionsOracleRepository) Get(userId string) ([]models.UserPermission, error) {
	rows, err := por.tx.Query(GET_PERMISSIONS_QUERY, userId)
	if err != nil {
		return []models.UserPermission{}, errors.New("Error making query: " + err.Error())
	}

	permissions := RowsToPermission(rows, false)
	return permissions, nil
}

func (por *PermissionsOracleRepository) Add(userPermission models.UserPermission) error{
	por.Delete(userPermission.UserId, userPermission.ServiceId)

	_, err := por.tx.Exec(CREATE_USER_PERMISSION_QUERY,
		userPermission.UserId, userPermission.ServiceId,
	)
	if err != nil {
		return errors.New("Permission could not be added")
	}

	return nil
}
func (por *PermissionsOracleRepository) Delete(userId string, serviceId string) error {
	_, err := por.tx.Exec(DELETE_PERMISSION_QUERY,
		userId, serviceId,
	)
	
	if err != nil {
		return errors.New("Permission could not be removed")
	}

	return nil
}
func (por *PermissionsOracleRepository) DeleteAll(userId string) error {
	_, err := por.tx.Exec(DELETE_ALL_USER_PERMISSIONS_QUERY,
		userId,
	)
	
	if err != nil {
		return errors.New("Permission could not be removed")
	}

	return nil
}

func (por *PermissionsOracleRepository) Update(userId string, userPermission []models.UserPermission) error {
	err := por.DeleteAll(userId)
	if err != nil {
		return errors.New("Error updating user permissions")
	}
	
	for _,v := range userPermission {
		err = por.Add(v)
		fmt.Println("HERE")
		if err != nil {
			return errors.New("Error updating user permissions")
		}
	}

	return nil
}


func RowsToPermission(rows *sql.Rows, containsPagination bool) []models.UserPermission {
	var permissions []models.UserPermission
	for rows.Next() {
		var perm models.UserPermission
		// var id string
		var r int
		if containsPagination {
			rows.Scan(&perm.UserId, &perm.ServiceId, &r)
		} else {
			rows.Scan(&perm.UserId, &perm.ServiceId)
		}

		permissions = append(permissions, perm)
	}

	defer rows.Close()

	return permissions
}