package models

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"github.com/Glintt/gAPI/api/user_permission"
	"github.com/Glintt/gAPI/api/utils"
)

type User struct {
	Id       bson.ObjectId `bson:"_id" json:"Id"`
	Username string
	Password string `json:",omitempty"`
	Email    string
	IsAdmin  bool
	ClientId string
	ClientSecret string
}

// GetInternalAPIUser Returns an internal user with admin permissions
func GetInternalAPIUser() User {
	return User{
		IsAdmin: true,
	}
}

// GeneratePassword encrypts password
func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// HasPermissionToAccessService check if user has permission to access service 
func (u *User) HasPermissionToAccessService(serviceID string) bool {
	if u.IsAdmin {
		return true
	}
	hasPermission , err := user_permission.UserHasPermissionToAccessService(u.Id.Hex(), serviceID)
	if err != nil {
		utils.LogMessage("Error getting user permissions: " + err.Error(), utils.DebugLogType)
		return false
	}

	return hasPermission
}