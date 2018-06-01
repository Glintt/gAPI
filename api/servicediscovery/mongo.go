package servicediscovery

import (
	"gAPIManagement/api/utils"
	"errors"
	
	"os"
	"regexp"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

const (
	SERVICES_COLLECTION = "services"
	SERVICE_GROUP_COLLECTION = "services_groups"
)

var MONGO_HOST string
var MONGO_DB string

func InitMongo() {
	MONGO_HOST = os.Getenv("MONGO_HOST")
	MONGO_DB = os.Getenv("MONGO_DB")
}

func ConnectToMongo() {
	session, err := mgo.Dial(MONGO_HOST)

	if err != nil {
		utils.LogMessage("error connecting to mongo")
	}

	db = session.DB(MONGO_DB)
}

func UpdateMongo(service Service, serviceExists Service) (string, int) {
	ConnectToMongo()

	err := db.C(SERVICES_COLLECTION).UpdateId(service.Id, &service)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service updated successfuly."}`, 201
}

func CreateServiceMongo(s Service) (string, int) {
	ConnectToMongo()

	s.Id = bson.NewObjectId()

	err := db.C(SERVICES_COLLECTION).Insert(&s)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service created successfuly."}`, 201
}

func ListServicesMongo(page int, filterQuery string) []Service {
	ConnectToMongo()

	skips := PAGE_LENGTH * (page - 1)

	var services []Service
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
	return services
}

func DeleteServiceMongo(s Service) (string, int) {
	service, err := FindMongo(s)

	if err != nil {
		return `{"error": true, "msg": "Not found"}`, 404
	}

	err = db.C(SERVICES_COLLECTION).Remove(&service)

	if err == nil {
		return `{"error": false, "msg": "Removed successfully."}`, 200
	}
	return `{"error": true, "msg": "Not found"}`, 404
}

func FindMongo(s Service) (Service, error) {
	ConnectToMongo()

	var services []Service

	f := func(c rune) bool {
		return c == '/'
	}
	uriParts := strings.FieldsFunc(s.MatchingURI, f)

	if s.Id == "" {
		s.Id = bson.NewObjectId()
	}
	query := bson.M{"$or": []bson.M{bson.M{"matchinguri": bson.RegEx{"/" + uriParts[0] + ".*", "i"}},bson.M{"_id": s.Id}}}
	db.C(SERVICES_COLLECTION).Find(query).All(&services)
	
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

func NormalizeServicesMongo() error {
	ConnectToMongo()

	var services []Service
	db.C(SERVICES_COLLECTION).Find(bson.M{}).All(&services)

	for _, rs := range services {
		rs.NormalizeService()

		db.C(SERVICES_COLLECTION).UpdateId(rs.Id, &rs)
	}	
	return nil
}