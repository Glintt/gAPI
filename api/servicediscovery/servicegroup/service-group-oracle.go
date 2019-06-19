package servicegroup

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/utils"

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

type ServiceGroupOracleRepository struct {
	Db      *sql.DB
	DbError error
	Tx      *sql.Tx
}

func (agmr *ServiceGroupOracleRepository) OpenTransaction() error {
	tx, err := agmr.Db.Begin()
	agmr.Tx = tx
	return err
}

func (agmr *ServiceGroupOracleRepository) CommitTransaction() {
	agmr.Tx.Commit()
}

func (agmr *ServiceGroupOracleRepository) RollbackTransaction() {
	agmr.Tx.Rollback()
}

func (agmr *ServiceGroupOracleRepository) Release() {
	database.CloseOracleConnection(agmr.Db)
}

// GetServiceGroups get list of service groups
func (sgor *ServiceGroupOracleRepository) GetServiceGroups() ([]ServiceGroup, error) {
	rows, err := sgor.Tx.Query(LIST_SERVICE_GROUP_V2)
	if err != nil {
		return []ServiceGroup{}, err
	}

	return RowsToServiceGroup(rows, false), nil
}

// GetServiceGroupById get service groups by id
func (sgor *ServiceGroupOracleRepository) GetServiceGroupById(serviceGroup string) (ServiceGroup, error) {
	rows, err := sgor.Tx.Query(GET_SERVICE_GROUP_BY_ID_OR_NAME, serviceGroup, serviceGroup)
	if err != nil {
		return ServiceGroup{}, err
	}

	groups := RowsToServiceGroup(rows, false)
	if len(groups) == 0 {
		return ServiceGroup{}, errors.New("Service group not found")
	}

	return groups[0], nil
}

// AddServiceToGroup add new service to a group
func (sgor *ServiceGroupOracleRepository) AddServiceToGroup(serviceGroupID string, serviceID string) error {
	_, err := sgor.Tx.Exec(ADD_SERVICE_TO_GROUP,
		serviceGroupID, serviceID,
	)
	return err
}

// RemoveServiceFromGroup remove a service from a group
func (sgor *ServiceGroupOracleRepository) RemoveServiceFromGroup(serviceGroupID string, serviceID string) error {
	_, err := sgor.Tx.Exec(REMOVE_SERVICE_FROM_GROUP,
		serviceID,
	)
	return err
}

// CreateServiceGroup create a new service group
func (sgor *ServiceGroupOracleRepository) CreateServiceGroup(serviceGroup ServiceGroup) error {
	_, err := sgor.Tx.Exec(CREATE_SERVICE_GROUP,
		serviceGroup.Id.Hex(), serviceGroup.Name, utils.BoolToInt(serviceGroup.IsReachable),
	)
	return err
}

// UpdateServiceGroup update an already existing service group
func (sgor *ServiceGroupOracleRepository) UpdateServiceGroup(serviceGroupID string, serviceGroup ServiceGroup) error {
	_, err := sgor.Tx.Exec(UPDATE_SERVICE_GROUP,
		serviceGroup.Name, utils.BoolToInt(serviceGroup.IsReachable), serviceGroupID,
	)
	return err
}

// DeleteServiceGroup delete a service group
func (sgor *ServiceGroupOracleRepository) DeleteServiceGroup(serviceGroupID string) error {
	_, err := sgor.Db.Exec(DELETE_SERVICE_GROUP,
		serviceGroupID,
	)
	return err
}

// RowsToServiceGroup convert sql rows to service group list
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
