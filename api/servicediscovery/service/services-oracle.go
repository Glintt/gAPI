package service

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Glintt/gAPI/api/database"
	userModels "github.com/Glintt/gAPI/api/users/models"
	"github.com/Glintt/gAPI/api/utils"

	//"github.com/Glintt/gAPI/api/user_permission"
	"strings"

	_ "gopkg.in/goracle.v2"
	"gopkg.in/mgo.v2/bson"
)

var SERVICE_COLUMNS = ` a.id, a.identifier,
	a.name, a.matchinguri, a.matchinguriregex, a.touri, a.protected, a.apidocumentation, a.iscachingactive, a.isactive, 
	a.healthcheckurl, a.lastactivetime, a.ratelimit, a.ratelimitexpirationtime, a.isreachable, 
	a.groupid, a.usegroupattributes, 
	a.servicemanagementhost, a.servicemanagementport,
	a.servicemanagementendpoints, a.hosts, a.protectedexclude, b.isreachable as groupreachable, a.oauth_clients_enabled as oauth_clients_enabled`

var LIST_SERVICES_ORACLE = `select ` + SERVICE_COLUMNS + ` 
	from gapi_services a left join gapi_services_groups b on a.groupid = b.id where 
	(upper(a.name) like upper(:name) or upper(a.matchinguri) like upper(:matchinguri)) `

var INSERT_SERVICE_ORACLE = `INSERT INTO gapi_services 
(
	id, identifier,name, matchinguri,matchinguriregex,touri,
	protected,apidocumentation,iscachingactive,isactive,healthcheckurl, lastactivetime, ratelimit, ratelimitexpirationtime, 
	isreachable, groupid, usegroupattributes, 
	servicemanagementhost, servicemanagementport,
	servicemanagementendpoints, hosts, protectedexclude
) 
VALUES(:id, :identifier,:name, :matchinguri,:matchinguriregex,:touri,
	:protected,:apidocumentation,:iscachingactive,:isactive,:healthcheckurl, :lastactivetime, :ratelimit, :ratelimitexpirationtime, 
	:isreachable, :groupid, :usegroupattributes, 
	:servicemanagementhost, :servicemanagementport,
	:servicemanagementendpoints, :hosts, :protectedexclude)`

var FIND_SERVICES_ORACLE = `select ` + SERVICE_COLUMNS +
	` from gapi_services a left join gapi_services_groups b on a.groupid = b.id where (a.id = :id or regexp_like(matchinguri, :matchinguriregex) or a.identifier = :identifier) `

var FIND_COLLISIONS = `select count(*) total 
from gapi_services a where a.matchinguri like :newmatchinguri and a.matchinguri <> :matchinguri`

var DELETE_SERVICES_ORACLE = `delete from gapi_services where id = :id`

var UPDATE_SERVICE_ORACLE = `UPDATE gapi_services
SET 
	identifier = :identifier,	
	name = :name, 
	matchinguri = :matchinguri, 
	matchinguriregex = :matchinguriregex, 
	touri = :touri, 
	protected = :protected, 
	apidocumentation = :apidocumentation, 
	iscachingactive = :iscachingactive, 
	isactive = :isactive, 	
	healthcheckurl = :healthcheckurl, 
	lastactivetime = :lastactivetime, 
	ratelimit = :ratelimit, 
	ratelimitexpirationtime = :ratelimitexpirationtime, 
	isreachable = :isreachable, 
	groupid = :groupid,
	usegroupattributes = :usegroupattributes, 
	servicemanagementhost = :servicemanagementhost, 
	servicemanagementport = :servicemanagementport,
	servicemanagementendpoints = :servicemanagementendpoints, 
	hosts = :hosts, 
	protectedexclude = :protectedexclude
WHERE id = :id`

var SERVICE_DISTINCT_HOSTS_ORACLE = `SELECT distinct domain from gapi_services_hosts`

var DELETE_SERVICE_HOSTS_ORACLE = "delete from gapi_services_hosts where service_id = :sid"
var INSERT_SERVICE_HOSTS_ORACLE = "INSERT INTO gapi_services_hosts(service_id, domain) VALUES(:sid, :domain)"

type ServiceOracleRepository struct {
	User userModels.User
}

const USER_PERMISSION_CHECK = "id in (select service_id from gapi_user_services_permissions c where c.user_id = '##USER_ID##') or 1 = ##IS_USER_ADMIN##"

func AppendPermissionFilterToQuery(query string, sTable string, sGTable string, user userModels.User) string {
	query = query + " and "
	permissionQuery := USER_PERMISSION_CHECK

	permissionQuery = strings.Replace(permissionQuery, "##USER_ID##", user.Id.Hex(), -1)
	isAdminValue := "0"
	if user.IsAdmin {
		isAdminValue = "1"
	}
	permissionQuery = strings.Replace(permissionQuery, "##IS_USER_ADMIN##", isAdminValue, -1)

	query = query + "(" + sTable + "." + permissionQuery + " or (" + sTable + ".isreachable = 1 or (" + sGTable + ".isreachable = 1 and " + sTable + ".usegroupattributes = 1))" + ")"
	return query
}

func (smo *ServiceOracleRepository) Update(service Service, serviceExists Service) (int, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return 400, err
	}

	service.NormalizeService()

	mgnendpoints, _ := json.Marshal(service.ServiceManagementEndpoints)
	hosts, _ := json.Marshal(service.Hosts)
	protectedexclude, _ := json.Marshal(service.ProtectedExclude)
	tx, err := db.Begin()
	if err != nil {
		database.CloseOracleConnection(db)
		return 400, err
	}

	err = smo.VerifyServiceMatchingCollision(service, tx)
	if err != nil {
		tx.Rollback()
		database.CloseOracleConnection(db)
		return 400, err
	}

	_, err = tx.Exec(UPDATE_SERVICE_ORACLE,
		service.GenerateIdentifier(),
		service.Name, service.MatchingURI, service.MatchingURIRegex, service.ToURI,
		utils.BoolToInt(service.Protected), service.APIDocumentation, utils.BoolToInt(service.IsCachingActive), utils.BoolToInt(service.IsActive),
		service.HealthcheckUrl, service.LastActiveTime, service.RateLimit, service.RateLimitExpirationTime, utils.BoolToInt(service.IsReachable),
		service.GroupId.Hex(), utils.BoolToInt(service.UseGroupAttributes),
		service.ServiceManagementHost, service.ServiceManagementPort,
		string(mgnendpoints), string(hosts), string(protectedexclude),

		service.Id.Hex(),
	)

	err = smo.DeleteHostsFromService(service, tx)
	if err != nil {
		tx.Rollback()
		database.CloseOracleConnection(db)
		return 400, err
	}
	err = smo.AddHostsToService(service, tx)
	if err != nil {
		tx.Rollback()
		database.CloseOracleConnection(db)
		return 400, err
	}

	tx.Commit()
	database.CloseOracleConnection(db)
	return 201, nil
}

func (smo *ServiceOracleRepository) AddHostsToService(s Service, tx *sql.Tx) error {
	for _, h := range s.Hosts {
		_, err := tx.Exec(INSERT_SERVICE_HOSTS_ORACLE, s.Id.Hex(), h)
		if err != nil {
			return err
		}
	}
	return nil
}

func (smo *ServiceOracleRepository) DeleteHostsFromService(s Service, tx *sql.Tx) error {
	_, err := tx.Exec(DELETE_SERVICE_HOSTS_ORACLE, s.Id.Hex())
	if err != nil {
		return err
	}
	return nil
}

func (smo *ServiceOracleRepository) VerifyServiceMatchingCollision(s Service, tx *sql.Tx) error {
	res, _ := tx.Query(FIND_COLLISIONS, s.MatchingURI+"%", s.MatchingURI)
	for res.Next() {
		var counter int
		res.Scan(&counter)

		if counter > 0 {
			return errors.New("Matching URI already exists for another service.")
		}
	}
	return nil
}

func (smo *ServiceOracleRepository) CreateService(s Service) (Service, error) {

	if s.ServiceManagementEndpoints == nil {
		s.ServiceManagementEndpoints = make(map[string]string)
	}
	if s.ProtectedExclude == nil {
		s.ProtectedExclude = make(map[string]string)
	}
	mgnendpoints, _ := json.Marshal(s.ServiceManagementEndpoints)
	hosts, _ := json.Marshal(s.Hosts)
	protectedexclude, _ := json.Marshal(s.ProtectedExclude)
	s.Id = bson.NewObjectId()

	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return Service{}, err
	}

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		database.CloseOracleConnection(db)
		return Service{}, err
	}

	err = smo.VerifyServiceMatchingCollision(s, tx)
	if err != nil {
		tx.Rollback()
		database.CloseOracleConnection(db)
		return Service{}, err
	}

	_, err = tx.Exec(INSERT_SERVICE_ORACLE,
		s.Id.Hex(), s.GenerateIdentifier(), s.Name, s.MatchingURI, s.MatchingURIRegex, s.ToURI,
		utils.BoolToInt(s.Protected), s.APIDocumentation, utils.BoolToInt(s.IsCachingActive), utils.BoolToInt(s.IsActive),
		s.HealthcheckUrl, s.LastActiveTime, s.RateLimit, s.RateLimitExpirationTime, utils.BoolToInt(s.IsReachable),
		s.GroupId.Hex(),
		utils.BoolToInt(s.UseGroupAttributes),
		s.ServiceManagementHost, s.ServiceManagementPort,
		string(mgnendpoints), string(hosts), string(protectedexclude),
	)

	err = smo.AddHostsToService(s, tx)

	if err != nil {
		tx.Rollback()
		database.CloseOracleConnection(db)
		return Service{}, err
	}

	tx.Commit()
	database.CloseOracleConnection(db)
	return s, nil
}

func (smo *ServiceOracleRepository) ListServices(page int, filterQuery string) []Service {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return nil
	}

	query := AppendPermissionFilterToQuery(LIST_SERVICES_ORACLE, "a", "b", smo.User)

	query = query + " order by a.id"
	var rows *sql.Rows
	var pagination = false
	if page >= 0 {
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

			rows, err = db.Query(query, "%"+filterQuery+"%", "%"+filterQuery+"%", page)
		pagination = true
	} else {
		rows, err = db.Query(query, "%"+filterQuery+"%", "%"+filterQuery+"%")
	}

	if err != nil {
		utils.LogMessage("Error running query", utils.DebugLogType)
		//defer rows.Close()
		database.CloseOracleConnection(db)
		return []Service{}
	}

	services := RowsToService(rows, pagination)
	database.CloseOracleConnection(db)
	return services
}

func (smo *ServiceOracleRepository) DeleteService(s Service) error {
	service, err := smo.Find(s)
	if err != nil {
		return err
	}

	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)

	tx, err := db.Begin()

	smo.DeleteHostsFromService(s, tx)

	_, err = tx.Exec(DELETE_SERVICES_ORACLE,
		service.Id.Hex(),
	)

	if err != nil {
		tx.Rollback()
		database.CloseOracleConnection(db)
		return err
	}

	tx.Commit()
	database.CloseOracleConnection(db)
	return nil
}

func (smo *ServiceOracleRepository) Find(s Service) (Service, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return Service{}, err
	}

	f := func(c rune) bool {
		return c == '/'
	}
	uriParts := strings.FieldsFunc(s.MatchingURI, f)

	if s.Id == "" {
		s.Id = bson.NewObjectId()
	}
	if len(uriParts) == 0 {
		uriParts = append(uriParts, "")
	}

	query := AppendPermissionFilterToQuery(FIND_SERVICES_ORACLE, "a", "b", smo.User)

	utils.LogMessage("Query = "+query, utils.DebugLogType)

	rows, err := db.Query(query, s.Id.Hex(),
		"/"+uriParts[0]+".*",
		s.Identifier)
	if err != nil {
		database.CloseOracleConnection(db)
		return Service{}, err
	}

	services := RowsToService(rows, false)

	database.CloseOracleConnection(db)
	return FindServiceInList(s, services)
}

func (smo *ServiceOracleRepository) NormalizeServices() error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	rows, err := db.Query(LIST_SERVICES_ORACLE, "%", "%")
	if err != nil {
		utils.LogMessage("Error running query", utils.DebugLogType)
		defer rows.Close()
		database.CloseOracleConnection(db)
		return err
	}

	services := RowsToService(rows, false)

	database.CloseOracleConnection(db)

	for _, rs := range services {
		rs.NormalizeService()

		go smo.Update(rs, rs)
	}

	return nil
}

func RowsToService(rows *sql.Rows, containsPagination bool) []Service {
	var services []Service
	for rows.Next() {
		var s Service
		var id, groupid string
		var mngendpoints,
			hosts,
			protectedexclude []byte

		var a int
		if containsPagination {
			rows.Scan(&id, &s.Identifier, &s.Name, &s.MatchingURI, &s.MatchingURIRegex, &s.ToURI, &s.Protected, &s.APIDocumentation, &s.IsCachingActive, &s.IsActive,
				&s.HealthcheckUrl, &s.LastActiveTime, &s.RateLimit, &s.RateLimitExpirationTime, &s.IsReachable,
				&groupid, &s.UseGroupAttributes,
				&s.ServiceManagementHost, &s.ServiceManagementPort,
				&mngendpoints, &hosts, &protectedexclude,
				&s.GroupVisibility, &s.OAuthClientsEnabled, &a,
			)
		} else {
			rows.Scan(&id, &s.Identifier, &s.Name, &s.MatchingURI, &s.MatchingURIRegex, &s.ToURI, &s.Protected, &s.APIDocumentation, &s.IsCachingActive, &s.IsActive,
				&s.HealthcheckUrl, &s.LastActiveTime, &s.RateLimit, &s.RateLimitExpirationTime, &s.IsReachable,
				&groupid, &s.UseGroupAttributes,
				&s.ServiceManagementHost, &s.ServiceManagementPort,
				&mngendpoints, &hosts, &protectedexclude,
				&s.GroupVisibility, &s.OAuthClientsEnabled,
			)
		}

		if bson.IsObjectIdHex(id) {
			s.Id = bson.ObjectIdHex(id)
		} else {
			s.Id = bson.NewObjectId()
		}
		if bson.IsObjectIdHex(groupid) {
			s.GroupId = bson.ObjectIdHex(groupid)
		}

		json.Unmarshal([]byte(mngendpoints), &s.ServiceManagementEndpoints)
		json.Unmarshal([]byte(hosts), &s.Hosts)
		json.Unmarshal([]byte(protectedexclude), &s.ProtectedExclude)

		services = append(services, s)
	}

	defer rows.Close()

	return services
}

func (smo *ServiceOracleRepository) ListAllAvailableHosts() ([]string, error) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(SERVICE_DISTINCT_HOSTS_ORACLE)
	if err != nil {
		utils.LogMessage("Error running query", utils.DebugLogType)
		defer rows.Close()
		database.CloseOracleConnection(db)
		return nil, err
	}

	var hosts []string

	for rows.Next() {
		var host string
		rows.Scan(&host)
		hosts = append(hosts, host)
	}
	defer rows.Close()
	database.CloseOracleConnection(db)
	return hosts, nil
}
