package users

import (
	"fmt"
	"gAPIManagement/api/config"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `bson:"_id" json:"Id"`
	Username string
	Password string `json:",omitempty"`
	Email    string
	IsAdmin  bool
}

const (
	USERS_COLLECTION = "users"
	SERVICE_NAME     = "gapi_users"
	PAGE_LENGTH      = 10
)

var UsersList []User

func InitUsers() {
	UserMethods[config.GApiConfiguration.ServiceDiscovery.Type]["init"].(func())()

	err := CreateUser(User{Username: "admin", Email: "admin@gapi.com", Password: "admin", IsAdmin: true})
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

	return UserMethods[config.GApiConfiguration.ServiceDiscovery.Type]["create"].(func(User) error)(user)
}

func UpdateUser(user User) error {
	return UserMethods[config.GApiConfiguration.ServiceDiscovery.Type]["update"].(func(User) error)(user)
}

func FindUsersByUsernameOrEmail(q string, page int) []User {
	return UserMethods[config.GApiConfiguration.ServiceDiscovery.Type]["findbyuseroremail"].(func(string, int) []User)(q, page)
}

func GetUserByUsername(username string) []User {
	return UserMethods[config.GApiConfiguration.ServiceDiscovery.Type]["findbyusername"].(func(string) []User)(username)
}
