package servicegroup

import (
	"database/sql"
	"errors"
	"gAPIManagement/api/database"
	"gAPIManagement/api/utils"

	"gopkg.in/mgo.v2/bson"
)

const (
	LIST_SERVICE_GROUP              = `select id, name, isreachable from gapi_services_groups`
	GET_SERVICE_GROUP_BY_ID_OR_NAME = `select id, name, isreachable from gapi_services_groups where id = :id or name = :name`
	ADD_SERVICE_TO_GROUP            = `update gapi_services set groupid = :groupid where id = :id`
	CREATE_SERVICE_GROUP            = `insert into gapi_services_groups(id, name, isreachable) values (:id,:name,:isreachable)`
	UPDATE_SERVICE_GROUP            = `update gapi_services_groups set name = :name, isreachable = :isreachable where id = :id`
	DELETE_SERVICE_GROUP            = `delete from gapi_services_groups where id = :id`
)

func GetServiceGroupsOracle() ([]ServiceGroup, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []ServiceGroup{}, err
	}

	rows, err := db.Query(LIST_SERVICE_GROUP)
	if err != nil {
		return []ServiceGroup{}, err
	}

	groups := RowsToServiceGroup(rows, false)
	database.CloseOracleConnection(db)

	return groups, nil
}

func GetServiceGroupByIdOracle(serviceGroup string) (ServiceGroup, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return ServiceGroup{}, err
	}

	rows, err := db.Query(GET_SERVICE_GROUP_BY_ID_OR_NAME, serviceGroup, serviceGroup)
	if err != nil {
		return ServiceGroup{}, err
	}

	groups := RowsToServiceGroup(rows, false)
	database.CloseOracleConnection(db)

	if len(groups) == 0 {
		return ServiceGroup{}, errors.New("Service group not found.")
	}

	return groups[0], nil
}

func AddServiceToGroupOracle(serviceGroupId string, serviceId string) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(ADD_SERVICE_TO_GROUP,
		serviceGroupId, serviceId,
	)

	database.CloseOracleConnection(db)
	return err
}

func CreateServiceGroupOracle(serviceGroup ServiceGroup) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(CREATE_SERVICE_GROUP,
		serviceGroup.Id.Hex(), serviceGroup.Name, utils.BoolToInt(serviceGroup.IsReachable),
	)

	database.CloseOracleConnection(db)
	return err
}

func UpdateServiceGroupOracle(serviceGroupId string, serviceGroup ServiceGroup) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(UPDATE_SERVICE_GROUP,
		serviceGroup.Name, utils.BoolToInt(serviceGroup.IsReachable), serviceGroupId,
	)

	database.CloseOracleConnection(db)
	return err
}

func DeleteServiceGroupOracle(serviceGroupId string) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(DELETE_SERVICE_GROUP,
		serviceGroupId,
	)

	database.CloseOracleConnection(db)
	return err
}

func RowsToServiceGroup(rows *sql.Rows, containsPagination bool) []ServiceGroup {
	var serviceGroups []ServiceGroup
	for rows.Next() {
		var serviceG ServiceGroup
		var id string
		var a int
		if containsPagination {
			rows.Scan(&id, &serviceG.Name, &serviceG.IsReachable, &a)
		} else {
			rows.Scan(&id, &serviceG.Name, &serviceG.IsReachable)
		}

		if bson.IsObjectIdHex(id) {
			serviceG.Id = bson.ObjectIdHex(id)
		} else {
			serviceG.Id = bson.NewObjectId()
		}

		serviceGroups = append(serviceGroups, serviceG)
	}

	defer rows.Close()

	return serviceGroups
}
