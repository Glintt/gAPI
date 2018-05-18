package servicediscovery

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

const (
	COLLECTION = "services"
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
		fmt.Println("error connecting to mongo")
	}

	db = session.DB(MONGO_DB)
}

func UpdateMongo(service Service, serviceExists Service) (string, int) {
	ConnectToMongo()

	err := db.C(COLLECTION).UpdateId(service.ID, &service)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service updated successfuly."}`, 201
}

func CreateServiceMongo(s Service) (string, int) {
	ConnectToMongo()

	s.ID = bson.NewObjectId()

	err := db.C(COLLECTION).Insert(&s)

	if err != nil {
		return `{"error" : true, "msg": "` + err.Error() + `"}`, 400
	}
	return `{"error" : false, "msg": "Service created successfuly."}`, 201
}

func ListServicesMongo() []Service {
	ConnectToMongo()

	var services []Service
	db.C(COLLECTION).Find(bson.M{}).All(&services)

	return services
}

func DeleteServiceMongo(matchingURI string) (string, int) {
	service, err := FindMongo(matchingURI)

	if err != nil {
		return `{"error": true, "msg": "Not found"}`, 404
	}

	err = db.C(COLLECTION).Remove(&service)

	if err == nil {
		return `{"error": false, "msg": "Removed successfully."}`, 200
	}
	return `{"error": true, "msg": "Not found"}`, 404
}

func FindMongo(toMatchUri string) (Service, error) {
	ConnectToMongo()

	var services []Service
	//db.C(COLLECTION).Find(bson.M{"matchinguri": toMatchUri}).All(&services)

	f := func(c rune) bool {
		return c == '/'
	}
	uriParts := strings.FieldsFunc(toMatchUri, f)

	db.C(COLLECTION).Find(bson.M{"matchinguri": bson.RegEx{"/" + uriParts[0] + ".*", "i"}}).All(&services)
	//{ $substrBytes: [ "toMatchUri", 0, 2 ] }

	//db.users.findOne({"username" : {$regexs : ".*son.*"}});

	for _, rs := range services {
		if (rs.MatchingURIRegex == "") {
			rs.MatchingURIRegex = GetMatchingURIRegex(rs.MatchingURI)
		}
		re := regexp.MustCompile(rs.MatchingURIRegex)
		if re.MatchString(toMatchUri) {
			return rs, nil
		}
	}

	return Service{}, errors.New("Not found.")
}

func NormalizeServicesMongo() error {
	ConnectToMongo()

	var services []Service
	db.C(COLLECTION).Find(bson.M{}).All(&services)

	for _, rs := range services {
		rs.NormalizeService()

		db.C(COLLECTION).UpdateId(rs.ID, &rs)
	}	
	return nil
}