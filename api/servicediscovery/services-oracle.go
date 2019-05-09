package servicediscovery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gAPIManagement/api/database"
	"gAPIManagement/api/utils"
	"strings"

	_ "gopkg.in/goracle.v2"
	"gopkg.in/mgo.v2/bson"
)

var SERVICE_COLUMNS = ` a.id, a.identifier,
	a.name, a.matchinguri, a.matchinguriregex, a.touri, a.protected, a.apidocumentation, a.iscachingactive, a.isactive, 
	a.healthcheckurl, a.lastactivetime, a.ratelimit, a.ratelimitexpirationtime, a.isreachable, 
	a.groupid, a.usegroupattributes, 
	a.servicemanagementhost, a.servicemanagementport,
	a.servicemanagementendpoints, a.hosts, a.protectedexclude `

var LIST_SERVICES_ORACLE = `select ` + SERVICE_COLUMNS + ` 
	from gapi_services a`

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

var FIND_SERVICES_ORACLE = `select	 ` + SERVICE_COLUMNS +
	` from gapi_services a where a.id = :id or regexp_like(matchinguri, :matchinguriregex) or a.identifier = :identifier`

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

func UpdateOracle(service Service, serviceExists Service) (string, int) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}

	service.NormalizeService()

	mgnendpoints, _ := json.Marshal(service.ServiceManagementEndpoints)
	hosts, _ := json.Marshal(service.Hosts)
	protectedexclude, _ := json.Marshal(service.ProtectedExclude)

	_, err = db.Exec(UPDATE_SERVICE_ORACLE,
		service.GenerateIdentifier(),
		service.Name, service.MatchingURI, service.MatchingURIRegex, service.ToURI,
		utils.BoolToInt(service.Protected), service.APIDocumentation, utils.BoolToInt(service.IsCachingActive), utils.BoolToInt(service.IsActive),
		service.HealthcheckUrl, service.LastActiveTime, service.RateLimit, service.RateLimitExpirationTime, utils.BoolToInt(service.IsReachable),
		service.GroupId.Hex(), utils.BoolToInt(service.UseGroupAttributes),
		service.ServiceManagementHost, service.ServiceManagementPort,
		string(mgnendpoints), string(hosts), string(protectedexclude),

		service.Id.Hex(),
	)

	database.CloseOracleConnection(db)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service updated successfuly."}`, 201
}

// tx, err := db.Begin()
// fmt.Println("ere")
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println("ere")
// defer tx.Rollback()
// fmt.Println("ere")
// stmt, err := tx.Prepare(INSERT_SERVICE_ORACLE)
// fmt.Println("ere")
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println("ere")
// defer stmt.Close() // danger!
// _, err = stmt.Exec("123123123")
// if err != nil {
// 	fmt.Println(err)
// 	return "nao deu", 0
// }
// fmt.Println("ere")
// err = tx.Commit()
// fmt.Println("ere")
// if err != nil {
// 	fmt.Println("fa")
// 	log.Fatal(err)
// }

func CreateServiceOracle(s Service) (string, int) {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}

	if s.ServiceManagementEndpoints == nil {
		s.ServiceManagementEndpoints = make(map[string]string)
	}
	if s.ProtectedExclude == nil {
		s.ProtectedExclude = make(map[string]string)
	}
	mgnendpoints, _ := json.Marshal(s.ServiceManagementEndpoints)
	hosts, _ := json.Marshal(s.Hosts)
	protectedexclude, _ := json.Marshal(s.ProtectedExclude)

	_, err = db.Exec(INSERT_SERVICE_ORACLE,
		bson.NewObjectId().Hex(), s.GenerateIdentifier(), s.Name, s.MatchingURI, s.MatchingURIRegex, s.ToURI,
		utils.BoolToInt(s.Protected), s.APIDocumentation, utils.BoolToInt(s.IsCachingActive), utils.BoolToInt(s.IsActive),
		s.HealthcheckUrl, s.LastActiveTime, s.RateLimit, s.RateLimitExpirationTime, utils.BoolToInt(s.IsReachable),
		s.GroupId.Hex(),
		utils.BoolToInt(s.UseGroupAttributes),
		s.ServiceManagementHost, s.ServiceManagementPort,
		string(mgnendpoints), string(hosts), string(protectedexclude),
	)

	database.CloseOracleConnection(db)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service created successfuly."}`, 201
}

func ListServicesOracle(page int, filterQuery string, viewAllPermission bool) []Service {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return nil
	}

	rows, err := db.Query(LIST_SERVICES_ORACLE)
	if err != nil {
		fmt.Println("Error running query")
		defer rows.Close()
		database.CloseOracleConnection(db)
		return []Service{}
	}

	database.CloseOracleConnection(db)
	return RowsToService(rows)
}

func DeleteServiceOracle(s Service) (string, int) {
	service, err := FindOracle(s)
	if err != nil {
		return `{"error": true, "msg": "Not found"}`, 404
	}

	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)

	_, err = db.Exec(DELETE_SERVICES_ORACLE,
		service.Id.Hex(),
	)

	database.CloseOracleConnection(db)
	if err == nil {
		return `{"error": false, "msg": "Removed successfully."}`, 200
	}
	return `{"error": true, "msg": "Service could not be removed"}`, 404
}

func FindOracle(s Service) (Service, error) {
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

	rows, err := db.Query(FIND_SERVICES_ORACLE, s.Id.Hex(),
		"/"+uriParts[0]+".*",
		s.Identifier)
	if err != nil {
		database.CloseOracleConnection(db)
		return Service{}, err
	}

	services := RowsToService(rows)

	database.CloseOracleConnection(db)
	return FindServiceInList(s, services)
}

func NormalizeServicesOracle() error {
	return nil
}

func RowsToService(rows *sql.Rows) []Service {
	var services []Service
	for rows.Next() {
		var s Service
		var id string
		var mngendpoints,
			hosts,
			protectedexclude []byte

		rows.Scan(&id, &s.Identifier, &s.Name, &s.MatchingURI, &s.MatchingURIRegex, &s.ToURI, &s.Protected, &s.APIDocumentation, &s.IsCachingActive, &s.IsActive,
			&s.HealthcheckUrl, &s.LastActiveTime, &s.RateLimit, &s.RateLimitExpirationTime, &s.IsReachable,
			&s.GroupId, &s.UseGroupAttributes,
			&s.ServiceManagementHost, &s.ServiceManagementPort,
			&mngendpoints, &hosts, &protectedexclude,
		)

		fmt.Println("hosts")
		fmt.Println(string(hosts))

		if bson.IsObjectIdHex(id) {
			s.Id = bson.ObjectIdHex(id)
		} else {
			s.Id = bson.NewObjectId()
		}
		json.Unmarshal([]byte(mngendpoints), &s.ServiceManagementEndpoints)
		json.Unmarshal([]byte(hosts), &s.Hosts)
		json.Unmarshal([]byte(protectedexclude), &s.ProtectedExclude)

		services = append(services, s)
	}

	defer rows.Close()

	return services
}
