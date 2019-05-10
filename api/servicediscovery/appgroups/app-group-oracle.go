package appgroups

import (
	"database/sql"
	"gAPIManagement/api/database"
	"gAPIManagement/api/servicediscovery"

	"gopkg.in/mgo.v2/bson"
)

var INSERT_APPLICATION_GROUP = `insert into gapi_services_apps_groups(id, name) values (:id, :name)`
var LIST_APPLICATION_GROUP = `select id, name from gapi_services_apps_groups where name like :name`
var GET_APPLICATION_GROUP_BY_ID = `select id, name from gapi_services_apps_groups where id = :id`
var GET_SERVICES_FOR_APPLICATION_GROUP = `select ` + servicediscovery.SERVICE_COLUMNS + ` 
	from gapi_services a left join gapi_services_groups b on a.groupid = b.id, gapi_services_apps_groups c where a.applicationgroupid = c.id and c.id = :id`

var DELETE_APPLICATION_GROUP = `delete from gapi_services_apps_groups where id = :id`
var UPDATE_APPLICATION_GROUP = `update gapi_services_apps_groups set name = :name where id = :id`
var GET_APPLICATION_GROUP_FOR_SERVICE = `select a.id, a.name from gapi_services_apps_groups a,
gapi_services b
 where b.id = :id and b.applicationgroupid = a.id`

var ASSOCIATE_APPLICATION_TO_GROUP = `update gapi_services set applicationgroupid = :groupid where id = :id`

func CreateApplicationGroupOracle(bodyMap ApplicationGroup) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}
	tx, err := db.Begin()

	_, err = tx.Exec(INSERT_APPLICATION_GROUP,
		bson.NewObjectId().Hex(), bodyMap.Name,
	)

	// TODO: update all services
	tx.Commit()
	database.CloseOracleConnection(db)

	return err
}

func GetApplicationGroupsOracle(page int, nameFilter string) []ApplicationGroup {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []ApplicationGroup{}
	}

	query := LIST_APPLICATION_GROUP
	query = `SELECT * FROM
			(
				SELECT a.*, rownum r__
				FROM
				(
					` + query + `
				) a
				WHERE rownum < ((:page * 10) + 1 )
			)
			WHERE r__ >= (((:page-1) * 10) + 1)`
	rows, err := db.Query(query, "%"+nameFilter+"%", page)

	groups := RowsToAppGroup(rows, true)
	database.CloseOracleConnection(db)

	return groups
}

func GetApplicationGroupByIdOracle(appGroupId string) (ApplicationGroup, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return ApplicationGroup{}, err
	}

	rows, err := db.Query(GET_APPLICATION_GROUP_BY_ID, appGroupId)
	if err != nil {
		return ApplicationGroup{}, err
	}

	appGroups := RowsToAppGroup(rows, false)
	database.CloseOracleConnection(db)

	if len(appGroups) == 0 {
		return ApplicationGroup{}, err
	}
	return appGroups[0], err
}

func GetServicesForApplicationGroupOracle(appGroup ApplicationGroup) ([]servicediscovery.Service, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []servicediscovery.Service{}, err
	}

	rows, err := db.Query(GET_SERVICES_FOR_APPLICATION_GROUP, appGroup.Id.Hex())
	if err != nil {
		return []servicediscovery.Service{}, err
	}

	services := servicediscovery.RowsToService(rows, false)
	database.CloseOracleConnection(db)

	return services, err
}

func DeleteApplicationGroupOracle(appGroupId string) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}
	tx, err := db.Begin()

	_, err = tx.Exec(DELETE_APPLICATION_GROUP,
		appGroupId,
	)

	// TODO: update all services
	tx.Commit()
	database.CloseOracleConnection(db)

	return err
}

func UpdateApplicationGroupOracle(appGroupId string, newGroup ApplicationGroup) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}
	tx, err := db.Begin()

	_, err = tx.Exec(UPDATE_APPLICATION_GROUP,
		newGroup.Name, appGroupId,
	)

	// TODO: update all services
	tx.Commit()
	database.CloseOracleConnection(db)

	return err
}

func FindServiceApplicationGroupOracle(serviceId string) ApplicationGroup {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return ApplicationGroup{}
	}

	rows, err := db.Query(GET_APPLICATION_GROUP_FOR_SERVICE,
		serviceId,
	)
	if err != nil {
		return ApplicationGroup{}
	}

	appGroups := RowsToAppGroup(rows, false)

	database.CloseOracleConnection(db)

	if len(appGroups) == 0 {
		return ApplicationGroup{}
	}
	return appGroups[0]
}

func RowsToAppGroup(rows *sql.Rows, containsPagination bool) []ApplicationGroup {
	var appGroups []ApplicationGroup
	for rows.Next() {
		var appG ApplicationGroup
		var id string
		var a int
		if containsPagination {
			rows.Scan(&id, &appG.Name, &a)
		} else {
			rows.Scan(&id, &appG.Name)
		}

		if bson.IsObjectIdHex(id) {
			appG.Id = bson.ObjectIdHex(id)
		} else {
			appG.Id = bson.NewObjectId()
		}
		appGroups = append(appGroups, appG)
	}

	defer rows.Close()

	return appGroups
}

func AddServiceToGroupOracle(appGroupId string, serviceId string) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(ASSOCIATE_APPLICATION_TO_GROUP,
		appGroupId, serviceId,
	)

	database.CloseOracleConnection(db)
	return err
}

func RemoveServiceFromGroupOracle(appGroupId string, serviceId string) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(ASSOCIATE_APPLICATION_TO_GROUP,
		"", serviceId,
	)

	database.CloseOracleConnection(db)
	return err
}
