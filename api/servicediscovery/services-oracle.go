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

var LIST_SERVICES_ORACLE = `select a.id, a.identifier,
	name, a.matchinguri, a.matchinguriregex, a.touri, a.protected, a.apidocumentation, a.iscachingactive, a.isactive, 
	a.healthcheckurl, a.lastactivetime, a.ratelimit, a.ratelimitexpirationtime, a.isreachable, a.usegroupattributes, 
	a.servicemanagementhost, a.servicemanagementport,
	a.servicemanagementendpoints, a.hosts, a.protectedexclude 	
	from gapi_services a`

var INSERT_SERVICE_ORACLE = `INSERT INTO gapi_services 
VALUES(:id, :identifier,:name, :matchinguri,:matchinguriregex,:touri,
	:protected,:apidocumentation,:iscachingactive,:isactive,:healthcheckurl, :lastactivetime, :ratelimit, :ratelimitexpirationtime, 
	:isreachable, :groupid, :usegroupattributes, 
	:servicemanagementhost, :servicemanagementport,
	:servicemanagementendpoints, :hosts, :protectedexclude)`

var FIND_SERVICES_ORACLE = `select a.id, a.identifier,
	name, a.matchinguri, a.matchinguriregex, a.touri, a.protected, a.apidocumentation, a.iscachingactive, a.isactive, 
	a.healthcheckurl, a.lastactivetime, a.ratelimit, a.ratelimitexpirationtime, a.isreachable, a.usegroupattributes, 
	a.servicemanagementhost, a.servicemanagementport,
	a.servicemanagementendpoints, a.hosts, a.protectedexclude 	
	from gapi_services a where a.id = :id or regexp_like(matchinguri, :matchinguriregex) or a.identifier = :identifier`

func UpdateOracle(service Service, serviceExists Service) (string, int) {
	return "", 0
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
		return "", 400
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

	fmt.Println(s)
	fmt.Println(string(mgnendpoints))
	fmt.Println(string(hosts))
	fmt.Println(string(protectedexclude))

	_, err = db.Exec(INSERT_SERVICE_ORACLE,
		bson.NewObjectId().Hex(), s.GenerateIdentifier(), s.Name, s.MatchingURI, s.MatchingURIRegex, s.ToURI,
		utils.BoolToInt(s.Protected), s.APIDocumentation, utils.BoolToInt(s.IsCachingActive), utils.BoolToInt(s.IsActive),
		s.HealthcheckUrl, s.LastActiveTime, s.RateLimit, s.RateLimitExpirationTime, utils.BoolToInt(s.IsReachable),
		s.GroupId.Hex(), utils.BoolToInt(s.UseGroupAttributes),
		s.ServiceManagementHost, s.ServiceManagementPort,
		string(mgnendpoints), string(hosts), string(protectedexclude),
	)

	if err != nil {
		fmt.Println(err)
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
	return "", 0
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

	fmt.Println(s.Identifier)
	fmt.Println(s.Identifier)
	fmt.Println("/" + uriParts[0] + ".*")
	rows, err := db.Query(FIND_SERVICES_ORACLE, s.Id.Hex(),
		"/"+uriParts[0]+".*",
		s.Identifier)
	if err != nil {
		fmt.Println(err)
		return Service{}, err
	}

	services := RowsToService(rows)

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
			&s.HealthcheckUrl, &s.LastActiveTime, &s.RateLimit, &s.RateLimitExpirationTime, &s.IsReachable, &s.UseGroupAttributes,
			&s.ServiceManagementHost, &s.ServiceManagementPort,
			&mngendpoints, &hosts, &protectedexclude,
		)

		if bson.IsObjectIdHex(id) {
			s.Id = bson.ObjectIdHex(id)
		} else {
			s.Id = bson.NewObjectId()
		}
		json.Unmarshal(mngendpoints, &s.ServiceManagementEndpoints)
		json.Unmarshal([]byte(hosts), &s.Hosts)
		json.Unmarshal([]byte(protectedexclude), &s.ProtectedExclude)

		services = append(services, s)
	}

	return services
}
