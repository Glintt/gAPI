package servicediscovery

import (
	"errors"
	"fmt"
	"os"

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

func ListServicesMongo(page int, filterQuery string) []Service {
	ConnectToMongo()

	skips := PAGE_LENGTH * (page - 1)

	var services []Service
	if page == -1 {
		db.C(COLLECTION).Find(bson.M{
			"$or": []bson.M{ 
				bson.M{"name": bson.RegEx{filterQuery+".*", ""}},
				bson.M{"matchinguri":bson.RegEx{filterQuery+".*", ""}}}}).Sort("matchinguri").All(&services)
	}else {
		db.C(COLLECTION).Find(bson.M{
			"$or": []bson.M{ 
				bson.M{"name": bson.RegEx{filterQuery+".*", ""}},
				bson.M{"matchinguri":bson.RegEx{filterQuery+".*", ""}}}}).Sort("matchinguri").Skip(skips).Limit(PAGE_LENGTH).All(&services)
	}
	return services
}

func DeleteServiceMongo(matchingURI string) (string, int) {
	service, err := FindMongo(GetMatchURI(matchingURI))

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
	db.C(COLLECTION).Find(bson.M{"matchinguri": toMatchUri}).All(&services)

	if len(services) > 0 {
		return services[0], nil
	}
	return Service{}, errors.New("Not found.")
}
