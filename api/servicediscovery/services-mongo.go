package servicediscovery

import (
	"gAPIManagement/api/database"
	"errors"
	"regexp"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

const (
	SERVICES_COLLECTION = "services"
	SERVICE_GROUP_COLLECTION = "services_groups"
)


func UpdateMongo(service Service, serviceExists Service) (string, int) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)
	
	err := db.C(SERVICES_COLLECTION).UpdateId(service.Id, &service)

	database.MongoDBPool.Close(session)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service updated successfuly."}`, 201
}

func CreateServiceMongo(s Service) (string, int) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	s.Id = bson.NewObjectId()
	err := db.C(SERVICES_COLLECTION).Insert(&s)

	database.MongoDBPool.Close(session)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service created successfuly."}`, 201
}

func ListServicesMongo(page int, filterQuery string) []Service {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var services []Service
	skips := PAGE_LENGTH * (page - 1)
	
	if page == -1 {
		db.C(SERVICES_COLLECTION).Find(bson.M{
			"$or": []bson.M{ 
				bson.M{"name": bson.RegEx{filterQuery+".*", ""}},
				bson.M{"matchinguri":bson.RegEx{filterQuery+".*", ""}}}}).Sort("matchinguri").All(&services)
	}else {
		db.C(SERVICES_COLLECTION).Find(bson.M{
			"$or": []bson.M{ 
				bson.M{"name": bson.RegEx{filterQuery+".*", ""}},
				bson.M{"matchinguri":bson.RegEx{filterQuery+".*", ""}}}}).Sort("matchinguri").Skip(skips).Limit(PAGE_LENGTH).All(&services)
	}

	database.MongoDBPool.Close(session)

	return services
}

func DeleteServiceMongo(s Service) (string, int) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	service, err := FindMongo(s)

	if err != nil {
		database.MongoDBPool.Close(session)
		
		return `{"error": true, "msg": "Not found"}`, 404
	}

	err = db.C(SERVICES_COLLECTION).Remove(&service)

	database.MongoDBPool.Close(session)

	if err == nil {
		return `{"error": false, "msg": "Removed successfully."}`, 200
	}
	return `{"error": true, "msg": "Not found"}`, 404
}

func FindMongo(s Service) (Service, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var services []Service

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
	query := bson.M{"$or": []bson.M{bson.M{"matchinguri": bson.RegEx{"/" + uriParts[0] + ".*", "i"}},bson.M{"_id": s.Id}}}
	db.C(SERVICES_COLLECTION).Find(query).All(&services)
	
	database.MongoDBPool.Close(session)

	for _, rs := range services {
		if (rs.MatchingURIRegex == "") {
			rs.MatchingURIRegex = GetMatchingURIRegex(rs.MatchingURI)
		}
		re := regexp.MustCompile(rs.MatchingURIRegex)
		if re.MatchString(s.MatchingURI) || rs.Id == s.Id {
			return rs, nil
		}
	}

	return Service{}, errors.New("Not found.")
}

func ListAllAvailableHosts() ([]string, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var hosts []string

	db.C(SERVICES_COLLECTION).Find(nil).Distinct("hosts", &hosts)
	
	database.MongoDBPool.Close(session)

	return hosts, nil
}

func NormalizeServicesMongo() error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var services []Service
	db.C(SERVICES_COLLECTION).Find(bson.M{}).All(&services)

	database.MongoDBPool.Close(session)
	
	for _, rs := range services {
		rs.NormalizeService()

		db.C(SERVICES_COLLECTION).UpdateId(rs.Id, &rs)
	}	
	return nil
}