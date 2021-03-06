package appgroups

import (
	"database/sql"
	"strings"

	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	"github.com/Glintt/gAPI/api/utils"

	"gopkg.in/mgo.v2/bson"
)

type AppGroupOracleRepository struct {
	Db      *sql.DB
	DbError error
	Tx      *sql.Tx
}

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

func (agmr *AppGroupOracleRepository) OpenTransaction() error {
	tx, err := agmr.Db.Begin()
	agmr.Tx = tx
	return err
}

func (agmr *AppGroupOracleRepository) CommitTransaction() {
	agmr.Tx.Commit()
}

func (agmr *AppGroupOracleRepository) RollbackTransaction() {
	agmr.Tx.Rollback()
}

func (agmr *AppGroupOracleRepository) Release() {
	database.CloseOracleConnection(agmr.Db)
}

// CreateApplicationGroup create application group
func (agmr *AppGroupOracleRepository) CreateApplicationGroup(bodyMap ApplicationGroup) error {
	id := bson.NewObjectId().Hex()
	_, err := agmr.Tx.Exec(INSERT_APPLICATION_GROUP,
		id, bodyMap.Name,
	)

	for _, rs := range bodyMap.Services {
		AddServiceToGroupQuery(id, rs.Hex(), agmr.Tx)
	}

	return err
}

// GetApplicationGroups get list of application groups
func (agmr *AppGroupOracleRepository) GetApplicationGroups(page int, nameFilter string) []ApplicationGroup {
	var groups []ApplicationGroup
	query := LIST_APPLICATION_GROUP_V2
	if page != -1 {
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
		rows, _ := agmr.Tx.Query(query, "%"+nameFilter+"%", page)
		groups = RowsToAppGroup(rows, true)
	} else {
		rows, _ := agmr.Tx.Query(query, "%"+nameFilter+"%")
		groups = RowsToAppGroup(rows, false)
	}

	return groups
}

// GetApplicationGroupByID get application group by id
func (agmr *AppGroupOracleRepository) GetApplicationGroupByID(appGroupID string) (ApplicationGroup, error) {
	rows, err := agmr.Tx.Query(GET_APPLICATION_GROUP_BY_ID, appGroupID)
	if err != nil {
		return ApplicationGroup{}, err
	}

	appGroups := RowsToAppGroup(rows, false)

	if len(appGroups) == 0 {
		return ApplicationGroup{}, err
	}
	return appGroups[0], err
}

// GetServicesForApplicationGroup get application group's services
func (agmr *AppGroupOracleRepository) GetServicesForApplicationGroup(appGroup ApplicationGroup) ([]service.Service, error) {
	rows, err := agmr.Tx.Query(GET_SERVICES_FOR_APPLICATION_GROUP, appGroup.Id.Hex())
	if err != nil {
		return []service.Service{}, err
	}

	return service.RowsToService(rows, false), err
}

// DeleteApplicationGroup delete application group by id
func (agmr *AppGroupOracleRepository) DeleteApplicationGroup(appGroupID string) error {
	_, err := agmr.Tx.Exec(DELETE_APPLICATION_GROUP,
		appGroupID,
	)

	// TODO: update all services
	return err
}

// UpdateApplicationGroup update application group with id appGroupID
func (agmr *AppGroupOracleRepository) UpdateApplicationGroup(appGroupID string, newGroup ApplicationGroup) error {
	_, err := agmr.Tx.Exec(UPDATE_APPLICATION_GROUP,
		newGroup.Name, appGroupID,
	)

	for _, rs := range newGroup.Services {
		AddServiceToGroupQuery(appGroupID, rs.Hex(), agmr.Tx)
	}
	return err
}

// FindServiceApplicationGroup get application group for service with id serviceID
func (agmr *AppGroupOracleRepository) FindServiceApplicationGroup(serviceID string) ApplicationGroup {
	rows, err := agmr.Tx.Query(GET_APPLICATION_GROUP_FOR_SERVICE,
		serviceID,
	)
	if err != nil {
		return ApplicationGroup{}
	}

	appGroups := RowsToAppGroup(rows, false)

	if len(appGroups) == 0 {
		return ApplicationGroup{}
	}
	return appGroups[0]
}

// AddServiceToGroup add service with id serviceID to application group with id appGroupID
func (agmr *AppGroupOracleRepository) AddServiceToGroup(appGroupID string, serviceID string) error {
	return AddServiceToGroupQuery(appGroupID, serviceID, agmr.Tx)
}

// AddServiceToGroupQuery associate service to application group
func AddServiceToGroupQuery(appGroupID string, serviceID string, tx *sql.Tx) error {
	_, err := tx.Exec(ASSOCIATE_APPLICATION_TO_GROUP,
		appGroupID, serviceID,
	)
	return err
}

// RemoveServiceFromGroup remove service from application group
func (agmr *AppGroupOracleRepository) RemoveServiceFromGroup(appGroupID string, serviceID string) error {
	_, err := agmr.Tx.Exec(ASSOCIATE_APPLICATION_TO_GROUP,
		"", serviceID,
	)
	return err
}

// UngroupedServices get list of services without any application group
func (agmr *AppGroupOracleRepository) UngroupedServices() []service.Service {
	rows, err := agmr.Tx.Query(UNGROUPED_SERVICES)
	if err != nil {
		utils.LogMessage(err.Error(), utils.ErrorLogType)
		return []service.Service{}
	}

	services := service.RowsToService(rows, false)
	return services
}

// RowsToAppGroup get applications groups from sql rows
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
