package servicediscovery

import (
	"encoding/json"
	"fmt"
	"gAPIManagement/api/database"

	_ "gopkg.in/goracle.v2"
	"gopkg.in/mgo.v2/bson"
)

var LIST_SERVICES_ORACLE = `select a.id, a.identifier,
	name, a.matchinguri, a.matchinguriregex, a.touri, a.protected, a.apidocumentation, a.iscachingactive, a.isactive, 
	a.healthcheckurl, a.lastactivetime, a.ratelimit, a.ratelimitexpirationtime, a.isreachable, a.usegroupattributes, 
	a.servicemanagementhost, a.servicemanagementport,
	a.servicemanagementendpoints, a.hosts, a.protectedexclude 	
	from gapi_services a`

func UpdateOracle(service Service, serviceExists Service) (string, int) {
	return "", 0
}

func CreateServiceOracle(s Service) (string, int) {
	return "", 0
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

		s.Id = bson.ObjectIdHex(id)
		json.Unmarshal(mngendpoints, &s.ServiceManagementEndpoints)
		json.Unmarshal([]byte(hosts), &s.Hosts)
		json.Unmarshal([]byte(protectedexclude), &s.ProtectedExclude)

		services = append(services, s)
	}

	database.CloseOracleConnection(db)
	return services
}

func DeleteServiceOracle(s Service) (string, int) {
	return "", 0
}

func FindOracle(s Service) (Service, error) {
	return Service{}, nil
}

func NormalizeServicesOracle() error {
	return nil
}
