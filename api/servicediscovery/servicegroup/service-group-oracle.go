package servicegroup

import (
	"database/sql"
	"errors"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/utils"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

const (
	LIST_SERVICE_GROUP    = `select id, name, isreachable from gapi_services_groups`
	LIST_SERVICE_GROUP_V2 = `select a.id, a.name, a.isreachable, listagg(b.id,', ') within group(order by b.id) services from gapi_services_groups a left join gapi_services b
on b.groupid = a.id  group by (a.id, a.name, a.isreachable )`

	GET_SERVICE_GROUP_BY_ID_OR_NAME = `select a.id, a.name, a.isreachable, listagg(b.id,', ') within group(order by b.id) services from gapi_services_groups a left join gapi_services b
on b.groupid = a.id  where a.id = :id or a.name = :name group by (a.id, a.name, a.isreachable )`
	ADD_SERVICE_TO_GROUP      = `update gapi_services set groupid = :groupid where id = :id`
	REMOVE_SERVICE_FROM_GROUP = `update gapi_services set groupid = '' where id = :id`
	CREATE_SERVICE_GROUP      = `insert into gapi_services_groups(id, name, isreachable) values (:id,:name,:isreachable)`
	UPDATE_SERVICE_GROUP      = `update gapi_services_groups set name = :name, isreachable = :isreachable where id = :id`
	DELETE_SERVICE_GROUP      = `delete from gapi_services_groups where id = :id`
)

func GetServiceGroupsOracle() ([]ServiceGroup, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []ServiceGroup{}, err
	}

	rows, err := db.Query(LIST_SERVICE_GROUP_V2)
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

func RemoveServiceFromGroupOracle(serviceGroupId string, serviceId string) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(REMOVE_SERVICE_FROM_GROUP,
		serviceId,
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
		var servicesList string
		if containsPagination {
			rows.Scan(&id, &serviceG.Name, &serviceG.IsReachable, &servicesList, &a)
		} else {
			rows.Scan(&id, &serviceG.Name, &serviceG.IsReachable, &servicesList)
		}

		servicesListArray := strings.Split(servicesList, ",")
		services := make([]bson.ObjectId, 0)
		for _, s := range servicesListArray {
			if !bson.IsObjectIdHex(strings.Trim(s, " ")) {
				continue
			}
			services = append(services, bson.ObjectIdHex(strings.Trim(s, " ")))
		}
		serviceG.Services = services
		serviceG.HostsEnabled = make([]string, 0)

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
