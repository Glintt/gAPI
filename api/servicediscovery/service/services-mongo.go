package service

import (
	"strings"

	"github.com/Glintt/gAPI/api/users"

	"errors"

	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"gopkg.in/mgo.v2/bson"
)

type ServiceMongoRepository struct {
	User users.User
}

func (smr *ServiceMongoRepository) Update(service Service, serviceExists Service) (int, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(constants.SERVICES_COLLECTION).UpdateId(service.Id, &service)

	database.MongoDBPool.Close(session)

	if err != nil {
		return 400, err
	}
	return 201, nil
}

func (smr *ServiceMongoRepository) CreateService(s Service) (Service, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	s.Id = bson.NewObjectId()
	err := db.C(constants.SERVICES_COLLECTION).Insert(&s)

	s.Identifier = s.GenerateIdentifier()

	database.MongoDBPool.Close(session)

	if err != nil {
		return Service{}, err
	}
	return s, nil
}

func (smr *ServiceMongoRepository) ListServices(page int, filterQuery string) []Service {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var services []Service
	skips := constants.PAGE_LENGTH * (page - 1)

	andQuery := []bson.M{
		bson.M{
			"$or": []bson.M{
				bson.M{"name": bson.RegEx{Pattern: filterQuery + ".*", Options: "i"}},
				bson.M{"matchinguri": bson.RegEx{Pattern: filterQuery + ".*", Options: "i"}}}},
	}

	if smr.User.Username == "" {
		visibilityQuery := bson.M{
			"$or": []bson.M{
				bson.M{"isreachable": true},
				bson.M{
					"$and": []bson.M{
						bson.M{"groupvisibility": true},
						bson.M{"usegroupattributes": true},
					},
				},
			}}
		andQuery = append(andQuery, visibilityQuery)
	}
	query := bson.M{
		"$and": andQuery,
	}

	if page == -1 {
		db.C(constants.SERVICES_COLLECTION).Find(query).Sort("matchinguri").All(&services)
	} else {
		db.C(constants.SERVICES_COLLECTION).Find(query).Sort("matchinguri").Skip(skips).Limit(constants.PAGE_LENGTH).All(&services)
	}

	database.MongoDBPool.Close(session)

	return services
}

func (smr *ServiceMongoRepository) DeleteService(s Service) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	service, err := smr.Find(s)

	if err != nil {
		database.MongoDBPool.Close(session)
		return errors.New("Service not found")
	}

	err = db.C(constants.SERVICES_COLLECTION).Remove(&service)

	database.MongoDBPool.Close(session)
	return err
}

func (smr *ServiceMongoRepository) Find(s Service) (Service, error) {
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
	query := bson.M{"$or": []bson.M{bson.M{"matchinguri": bson.RegEx{"/" + uriParts[0] + ".*", "i"}}, bson.M{"_id": s.Id}, bson.M{"identifier": s.Identifier}}}
	db.C(constants.SERVICES_COLLECTION).Find(query).All(&services)

	database.MongoDBPool.Close(session)

	return FindServiceInList(s, services)
}

func (smr *ServiceMongoRepository) ListAllAvailableHosts() ([]string, error) {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var hosts []string

	db.C(constants.SERVICES_COLLECTION).Find(nil).Distinct("hosts", &hosts)

	database.MongoDBPool.Close(session)

	return hosts, nil
}

func (smr *ServiceMongoRepository) NormalizeServices() error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var services []Service
	db.C(constants.SERVICES_COLLECTION).Find(bson.M{}).All(&services)

	for _, rs := range services {
		rs.NormalizeService()

		db.C(constants.SERVICES_COLLECTION).UpdateId(rs.Id, &rs)
	}

	database.MongoDBPool.Close(session)

	return nil
}
