package authentication

import (
	"errors"
	"encoding/json"
	"gAPIManagement/api/config"
	"io/ioutil"
)


type User struct {
	Id string
	Username string
	Password string
	Email string
}

var UsersList []User

func LoadUsers(){
	usersJson, err := ioutil.ReadFile(config.CONFIGS_LOCATION + "users.json")

	if err != nil {
		return
	}

	json.Unmarshal(usersJson, &UsersList)
	return
}

func FindUserByUsername(username string) (User, error) {
	for _, user := range UsersList {
		if user.Username == username{
			return user, nil
		}
	}

	return User{}, errors.New("User not found.")
}