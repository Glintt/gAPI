package users

import (
	"github.com/Glintt/gAPI/api/users/models"
	"github.com/Glintt/gAPI/api/users/providers"
	"github.com/Glintt/gAPI/api/utils"
	"gopkg.in/mgo.v2/bson"
	"errors"
)
const (
	SERVICE_NAME     = "gapi_users"
)

type UserService struct {
	User          models.User
	UserRepos 		providers.UserRepository
}

// InitUsers init users service
func InitUsers() {
	userService := UserService{User: models.GetInternalAPIUser()}
	userService.createRepository()
	userService.UserRepos.OpenTransaction()

	userService.UserRepos.InitUsers()
	
	err := userService.CreateUser(models.User{Username: "admin", Email: "admin@gapi.com", Password: "admin", IsAdmin: true})
	if err != nil {
		utils.LogMessage(err.Error(), utils.ErrorLogType)
	}
	releaseConnection(&userService)
}

// NewUserServiceWithUser create user service
func NewUserServiceWithUser(user models.User) (UserService, error) {
	userServ := UserService{User: user}
	err := userServ.createRepository()
	return userServ, err
}

func releaseConnection(us *UserService) {
	us.UserRepos.CommitTransaction()
	us.UserRepos.Release()	
}

func (us *UserService) createRepository() error {
	us.UserRepos = providers.NewUserRepository(us.User)
	if us.UserRepos == nil {
		return errors.New("Could not get application group repository")
	}
	us.UserRepos.OpenTransaction()
	return nil
}

// CreateUser creates a new user
func (us *UserService) CreateUser(user models.User) error {
	us.UserRepos.OpenTransaction()
	hashedPwd, _ := models.GeneratePassword(user.Password)
	user.Password = string(hashedPwd)
	user.Id = bson.NewObjectId()

	err := us.UserRepos.CreateUser(user)
	releaseConnection(us)
	return err
}

// UpdateUser update an existing user
func (us *UserService) UpdateUser(user models.User) error {
	us.UserRepos.OpenTransaction()
	err := us.UserRepos.UpdateUser(user)
	releaseConnection(us)
	return err
}
// FindUsersByUsernameOrEmail search user by email or username
func (us *UserService) FindUsersByUsernameOrEmail(q string, page int) []models.User {
	us.UserRepos.OpenTransaction()
	users := us.UserRepos.FindUsersByUsernameOrEmail(q, page)
	releaseConnection(us)
	return users
}
// GetUserByUsername search user by username
func (us *UserService) GetUserByUsername(username string) []models.User {
	us.UserRepos.OpenTransaction()
	users := us.UserRepos.GetUserByUsername(username)
	releaseConnection(us)
	return users
}
