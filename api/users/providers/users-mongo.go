package providers

import (
	"github.com/Glintt/gAPI/api/database"
	models "github.com/Glintt/gAPI/api/users/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	USERS_COLLECTION = "users"
)
type UserMongoRepository struct {
	Session    *mgo.Session
	Db         *mgo.Database
	Collection *mgo.Collection
}

// OpenTransaction opens a new database transaction
func (ur *UserMongoRepository) OpenTransaction() error {
	return nil
}
// CommitTransaction commits a database transaction
func (ur *UserMongoRepository) CommitTransaction() {
}
// RollbackTransaction rollbacks a database transaction
func (ur *UserMongoRepository) RollbackTransaction() {
}
// Release releases a database connection to the pool
func (ur *UserMongoRepository) Release() {
	database.MongoDBPool.Close(ur.Session)
}

// InitUsers inits user required modules
func (ur *UserMongoRepository) InitUsers() {
	if !database.IsConnectionDone {
		if err := database.InitDatabaseConnection(); err != nil {
			panic(err.Error())
		}
	}

	index := mgo.Index{
		Key:        []string{"username", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := ur.Collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// CreateUser create a new user
func (ur *UserMongoRepository) CreateUser(user models.User) error {
	hashedPwd, _ := models.GeneratePassword(user.Password)
	user.Password = string(hashedPwd)
	user.Id = bson.NewObjectId()

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(USERS_COLLECTION).Insert(&user)

	database.MongoDBPool.Close(session)

	return err
}

// UpdateUser update an existing user
func (ur *UserMongoRepository) UpdateUser(user models.User) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(USERS_COLLECTION).UpdateId(user.Id, &user)

	database.MongoDBPool.Close(session)

	return err
}

// FindUsersByUsernameOrEmail return a list of users searched by username or email
func (ur *UserMongoRepository) FindUsersByUsernameOrEmail(q string, page int) []models.User {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	p := ".*" + q + ".*"
	query := bson.M{"$or": []bson.M{bson.M{"username": bson.RegEx{p, "i"}},
		bson.M{"email": bson.RegEx{p, "i"}}}}

	skips := PAGE_LENGTH * (page - 1)
	var users []models.User

	db.C(USERS_COLLECTION).Find(query).Select(bson.M{"password": 0}).Sort("username").Skip(skips).Limit(PAGE_LENGTH).All(&users)

	database.MongoDBPool.Close(session)

	return users
}

// GetUserByUsername return a list of users search by username
func (ur *UserMongoRepository) GetUserByUsername(username string) []models.User {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	query := bson.M{"$or": []bson.M{bson.M{"username": username}}}

	var users []models.User

	db.C(USERS_COLLECTION).Find(query).All(&users)

	database.MongoDBPool.Close(session)

	return users
}

func (ur *UserMongoRepository) GetUserByIdentifier(id string) models.User {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	query := bson.M{"$or": []bson.M{bson.M{"id": bson.ObjectIdHex(id)}}}

	var user models.User

	db.C(USERS_COLLECTION).Find(query).One(&user)

	database.MongoDBPool.Close(session)

	return user
}
