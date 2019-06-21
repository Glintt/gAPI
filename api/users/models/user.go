package models

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"github.com/Glintt/gAPI/api/user_permission"
	"github.com/Glintt/gAPI/api/utils"
	oauthClientModels "github.com/Glintt/gAPI/api/oauth_clients/models"
	"github.com/Glintt/gAPI/api/oauth_clients"
)

type User struct {
	Id       bson.ObjectId `bson:"_id" json:"Id"`
	Username string
	Password string `json:",omitempty"`
	Email    string
	IsAdmin  bool
	OAuthClients []oauthClientModels.OAuthClient `json:oauth_clients""`
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
// GetOAuthClients Get list of oauth clients associated to user 
func (u *User) GetOAuthClients() []oauthClientModels.OAuthClient {
	return oauth_clients.FindForUser(u.Id.Hex())
}