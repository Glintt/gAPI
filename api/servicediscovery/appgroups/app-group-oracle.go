package appgroups

import (
	"database/sql"
	"fmt"
	"gAPIManagement/api/database"
	"gAPIManagement/api/servicediscovery/service"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

var INSERT_APPLICATION_GROUP = `insert into gapi_services_apps_groups(id, name) values (:id, :name)`
var LIST_APPLICATION_GROUP_V2 = `select a.id, a.name, listagg(b.id,', ') within group(order by b.id) services from gapi_services_apps_groups a left join gapi_services b
on b.applicationgroupid = a.id  where a.name like :name group by (a.id, a.name)`
var LIST_APPLICATION_GROUP = `select a.id, a.name from gapi_services_apps_groups a where a.name like :name`
var GET_APPLICATION_GROUP_BY_ID = `select id, name, '' as services from gapi_services_apps_groups where id = :id`
var GET_SERVICES_FOR_APPLICATION_GROUP = `select ` + service.SERVICE_COLUMNS + ` 
	from gapi_services a left join gapi_services_groups b on a.groupid = b.id, gapi_services_apps_groups c where a.applicationgroupid = c.id and c.id = :id`

var DELETE_APPLICATION_GROUP = `delete from gapi_services_apps_groups where id = :id`
var UPDATE_APPLICATION_GROUP = `update gapi_services_apps_groups set name = :name where id = :id`
var GET_APPLICATION_GROUP_FOR_SERVICE = `select a.id, a.name, '' as services from gapi_services_apps_groups a,
gapi_services b
 where b.id = :id and b.applicationgroupid = a.id`

var ASSOCIATE_APPLICATION_TO_GROUP = `update gapi_services set applicationgroupid = :groupid where id = :id`
var UNGROUPED_SERVICES = `select ` + service.SERVICE_COLUMNS + `  from gapi_services a left join gapi_services_groups b on a.groupid = b.id where applicationgroupid is null`

func CreateApplicationGroupOracle(bodyMap ApplicationGroup) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}
	tx, err := db.Begin()

	id := bson.NewObjectId().Hex()
	_, err = tx.Exec(INSERT_APPLICATION_GROUP,
		id, bodyMap.Name,
	)

	for _, rs := range bodyMap.Services {
		AddServiceToGroupQuery(id, rs.Hex(), tx)
	}

	tx.Commit()
	database.CloseOracleConnection(db)

	return err
}

func GetApplicationGroupsOracle(page int, nameFilter string) []ApplicationGroup {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []ApplicationGroup{}
	}

	query := LIST_APPLICATION_GROUP_V2
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

func GetServicesForApplicationGroupOracle(appGroup ApplicationGroup) ([]service.Service, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []service.Service{}, err
	}

	rows, err := db.Query(GET_SERVICES_FOR_APPLICATION_GROUP, appGroup.Id.Hex())
	if err != nil {
		return []service.Service{}, err
	}

	services := service.RowsToService(rows, false)
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

	for _, rs := range newGroup.Services {
		AddServiceToGroupQuery(appGroupId, rs.Hex(), tx)
	}

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
		var servicesList string
		var a int
		if containsPagination {
			rows.Scan(&id, &appG.Name, &servicesList, &a)
		} else {
			rows.Scan(&id, &appG.Name, &servicesList)
		}

		servicesListArray := strings.Split(servicesList, ",")
		if bson.IsObjectIdHex(id) {
			appG.Id = bson.ObjectIdHex(id)
		} else {
			appG.Id = bson.NewObjectId()
		}

		services := make([]bson.ObjectId, 0)
		for _, s := range servicesListArray {
			if !bson.IsObjectIdHex(strings.Trim(s, " ")) {
				continue
			}
			services = append(services, bson.ObjectIdHex(strings.Trim(s, " ")))
		}

		appG.Services = services
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

	tx, err := db.Begin()

	AddServiceToGroupQuery(appGroupId, serviceId, tx)

	tx.Commit()
	database.CloseOracleConnection(db)
	return err
}

func AddServiceToGroupQuery(appGroupId string, serviceId string, tx *sql.Tx) error {
	_, err := tx.Exec(ASSOCIATE_APPLICATION_TO_GROUP,
		appGroupId, serviceId,
	)
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

func UngroupedServicesOracle() []service.Service {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []service.Service{}
	}

	rows, err := db.Query(UNGROUPED_SERVICES)
	if err != nil {
		fmt.Println(err)
		return []service.Service{}
	}

	services := service.RowsToService(rows, false)

	database.CloseOracleConnection(db)
	return services
}
