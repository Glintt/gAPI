package users

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gAPIManagement/api/database"
)

type User struct {
	Id bson.ObjectId `bson:"_id" json:"Id"`
	Username string
	Password string `json:",omitempty"`
	Email string
	IsAdmin bool
}

const (
	USERS_COLLECTION = "users"
	SERVICE_NAME = "gapi_users"
	PAGE_LENGTH = 10
)

var UsersList []User

func InitUsers() {	
	if ! database.IsConnectionDone {
		if err := database.InitDatabaseConnection(); err != nil {
			panic(err.Error())
		}
	}
	
	_, db := database.GetSessionAndDB(database.MONGO_DB)
	
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
		panic(err)
	}

	err = CreateUser(User{Username: "admin", Email: "admin@gapi.com", Password: "admin", IsAdmin: true})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CreateUser(user User) error {
	hashedPwd, _ := GeneratePassword(user.Password)
	user.Password = string(hashedPwd)
	user.Id = bson.NewObjectId()

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(USERS_COLLECTION).Insert(&user)

	database.MongoDBPool.Close(session)

	return err
}

func UpdateUser(user User) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(USERS_COLLECTION).UpdateId(user.Id, &user)

	database.MongoDBPool.Close(session)

	return err
}

func FindUsersByUsernameOrEmail(q string, page int ) []User {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	p := ".*" + q + ".*"
	query := bson.M{"$or": []bson.M{bson.M{"username": bson.RegEx{p, "i"}}, 
									bson.M{"email": bson.RegEx{p ,  "i"}}}}

	skips := PAGE_LENGTH * (page - 1)
	var users []User
	
	db.C(USERS_COLLECTION).Find(query).Select(bson.M{"password": 0}).Sort("username").Skip(skips).Limit(PAGE_LENGTH).All(&users)
	
	database.MongoDBPool.Close(session)
	
	return users
}

func GetUserByUsername(username string) []User {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	query := bson.M{"$or": []bson.M{bson.M{"username": username}}}

	var users []User
	
	db.C(USERS_COLLECTION).Find(query).All(&users)
	
	database.MongoDBPool.Close(session)
	
	return users
}