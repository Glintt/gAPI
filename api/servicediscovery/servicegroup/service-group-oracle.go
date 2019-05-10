package servicegroup

import (
	"database/sql"
	"gAPIManagement/api/database"

	"gopkg.in/mgo.v2/bson"
)

const LIST_SERVICE_GROUP = `select id, name, isreachable from gapi_services_groups`

func GetServiceGroupsOracle() ([]ServiceGroup, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []ServiceGroup{}, err
	}

	rows, err := db.Query(LIST_SERVICE_GROUP)

	groups := RowsToServiceGroup(rows, false)
	database.CloseOracleConnection(db)

	return groups, nil
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
