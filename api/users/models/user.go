package models

import (
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
