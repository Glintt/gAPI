package users

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"os"
	"gopkg.in/mgo.v2/bson"
	"gAPIManagement/api/database"
)

type User struct {
	Id bson.ObjectId `bson:"_id" json:"Id"`
	Username string
	Password string `json:"-"`
	Email string
	IsAdmin bool
}

const (
	USERS_COLLECTION = "users"
	SERVICE_NAME = "gapi_users"
	PAGE_LENGTH = 10
)
var MONGO_DB = ""
var UsersList []User

func InitUsers() {	
	MONGO_DB = os.Getenv("MONGO_DB")
	
	_, db := database.GetSessionAndDB(MONGO_DB)
	
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

func CreateUser(user User) error {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPwd)
	user.Id = bson.NewObjectId()

	session, db := database.GetSessionAndDB(MONGO_DB)

	err := db.C(USERS_COLLECTION).Insert(&user)

	database.MongoDBPool.Close(session)

	return err
}

func FindUsersByUsernameOrEmail(q string, page int ) []User {
	session, db := database.GetSessionAndDB(MONGO_DB)

	query := bson.M{"$or": []bson.M{bson.M{"username": bson.RegEx{"/" + q + ".*", "i"}},bson.M{"email": q}}}

	skips := PAGE_LENGTH * (page - 1)
	var users []User
	
	db.C(USERS_COLLECTION).Find(query).Sort("matchinguri").Skip(skips).Limit(PAGE_LENGTH).All(&users)
	
	database.MongoDBPool.Close(session)
	
	return users
}

func GetUserByUsername(username string) []User {
	session, db := database.GetSessionAndDB(MONGO_DB)

	query := bson.M{"$or": []bson.M{bson.M{"username": username}}}

	var users []User
	
	db.C(USERS_COLLECTION).Find(query).All(&users)
	
	database.MongoDBPool.Close(session)
	
	return users
}