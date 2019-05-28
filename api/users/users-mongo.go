package users

import (
	"github.com/Glintt/gAPI/api/database"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InitUsersMongo() {
	if !database.IsConnectionDone {
		if err := database.InitDatabaseConnection(); err != nil {
			panic(err.Error())
		}
	}

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	userCollection := db.C(USERS_COLLECTION)

	index := mgo.Index{
		Key:        []string{"username", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := userCollection.EnsureIndex(index)
	if err != nil {
		database.MongoDBPool.Close(session)

		panic(err)
	}

}

func CreateUserMongo(user User) error {
	hashedPwd, _ := GeneratePassword(user.Password)
	user.Password = string(hashedPwd)
	user.Id = bson.NewObjectId()

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(USERS_COLLECTION).Insert(&user)

	database.MongoDBPool.Close(session)

	return err
}

func UpdateUserMongo(user User) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(USERS_COLLECTION).UpdateId(user.Id, &user)

	database.MongoDBPool.Close(session)

	return err
}

func FindUsersByUsernameOrEmailMongo(q string, page int) []User {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	p := ".*" + q + ".*"
	query := bson.M{"$or": []bson.M{bson.M{"username": bson.RegEx{p, "i"}},
		bson.M{"email": bson.RegEx{p, "i"}}}}

	skips := PAGE_LENGTH * (page - 1)
	var users []User

	db.C(USERS_COLLECTION).Find(query).Select(bson.M{"password": 0}).Sort("username").Skip(skips).Limit(PAGE_LENGTH).All(&users)

	database.MongoDBPool.Close(session)

	return users
}

func GetUserByUsernameMongo(username string) []User {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	query := bson.M{"$or": []bson.M{bson.M{"username": username}}}

	var users []User

	db.C(USERS_COLLECTION).Find(query).All(&users)

	database.MongoDBPool.Close(session)

	return users
}
