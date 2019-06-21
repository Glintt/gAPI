package controllers

import (
	"encoding/json"
	"github.com/Glintt/gAPI/api/http"
	userModels "github.com/Glintt/gAPI/api/users/models"
	"github.com/Glintt/gAPI/api/users"
	"github.com/Glintt/gAPI/api/utils"

	routing "github.com/qiangxue/fasthttp-routing"
)

// UsersServiceName returns the name of this controller service group
func UsersServiceName() string {
	return users.SERVICE_NAME
}

// GetUserHandler handle GET /user/<username>
func GetUserHandler(c *routing.Context) error {
	// Get user service
	userService := getUserService(c)

	user := userService.GetUserByUsername(c.Param("username"))
	if len(user) == 0 {
		return http.NotFound(c, "User not found", UsersServiceName())
	}

	user[0].Password = ""
	userJSON, _ := json.Marshal(user[0])
	return http.Ok(c, string(userJSON), UsersServiceName())
}

// FindUsersHandler handle GET /users
func FindUsersHandler(c *routing.Context) error {
	// Get page query parameter
	page, err := http.ParsePageParam(c)
	if err != nil {
		return http.Error(c, err.Error(), 400, UsersServiceName())
	}

	// Get search query parameter
	searchQuery := ""
	if c.QueryArgs().Has("q") {
		searchQuery = string(c.QueryArgs().Peek("q"))
	}

	// Get user service
	userService := getUserService(c)

	users := userService.FindUsersByUsernameOrEmail(searchQuery, page)

	for i := range users {
		users[i].OAuthClients = users[i].GetOAuthClients()
	}
	if users == nil {
		return http.Ok(c, `[]`, UsersServiceName())
	}

	userJSON, _ := json.Marshal(users)
	return http.Ok(c, string(userJSON), UsersServiceName())
}

// CreateUserHandler handle POST /users
func CreateUserHandler(c *routing.Context) error {
	// Try parse post body
	var user userModels.User
	err := json.Unmarshal(c.Request.Body(), &user)
	if err != nil {
		return http.Error(c, "User could not be created", 400, UsersServiceName())
	}

	// Get user service
	userService := getUserService(c)

	err = userService.CreateUser(user)

	if err != nil {
		return http.Error(c, "User could not be created", 400, UsersServiceName())
	} 
	return http.Created(c, "User created successfuly", UsersServiceName())
}

// UpdateUserByAdminHandler handle PUT /users
func UpdateUserByAdminHandler(c *routing.Context) error {
	// Get user service
	userService := getUserService(c)

	// Get user with provided username
	usersList := userService.GetUserByUsername(c.Param("username"))
	if len(usersList) == 0 {
		return http.Error(c, "User not found", 404, UsersServiceName())
	}

	user := usersList[0]

	// Update user information
	allowedToUpdate := []string{"Email", "Password", "IsAdmin"}
	user, err := updateUserUsingUpdateBody(user, c.Request.Body(), allowedToUpdate)
	if err != nil {
		return http.Error(c, "User could not be updated", 400, UsersServiceName())
	}


	err = userService.UpdateUser(user)

	if err != nil {
		return http.Error(c, "User could not be updated", 400, UsersServiceName())
	} 
	return http.Created(c, "User updated successfuly", UsersServiceName())
}

// UpdateUserHandler handle PUT /users
func UpdateUserHandler(c *routing.Context) error {
	// Get user service
	userService := getUserService(c)

	// Check if user to update is the same as the user that is calling the service
	requestUser := string(c.Request.Header.Peek("User"))
	if requestUser != c.Param("username") {
		return http.Error(c, "User could not be updated", 400, UsersServiceName())
	}

	// Get user by username
	usersList := userService.GetUserByUsername(c.Param("username"))
	if len(usersList) == 0 {
		return http.Error(c, "User not found", 404, UsersServiceName())
	}
	user := usersList[0]

	// Update user information
	allowedToUpdate := []string{"email", "password"}
	user, err := updateUserUsingUpdateBody(user, c.Request.Body(), allowedToUpdate)
	if err != nil {
		return http.Error(c, "User could not be updated", 400, UsersServiceName())		
	}
	
	err = userService.UpdateUser(user)

	if err != nil {
		return http.Error(c, "User could not be updated", 400, UsersServiceName())
	} 
	return http.Created(c, "User updated successfuly", UsersServiceName())
}

func updateUserUsingUpdateBody(user userModels.User, body []byte, allowedToUpdate []string) (userModels.User, error) {
	// Try parse body
	var userUpdateObj map[string]interface{}
	err := json.Unmarshal(body, &userUpdateObj)
	if err != nil {
		return user, err
	}

	if email, ok := userUpdateObj["Email"]; ok {
		user.Email = email.(string)
	}
	if pwd, ok := userUpdateObj["Password"]; ok {
		hashedPwd, _ := userModels.GeneratePassword(pwd.(string))
		user.Password = string(hashedPwd)
	}

	if utils.ArrayContainsString(allowedToUpdate, "IsAdmin") {
		if isAdmin, ok := userUpdateObj["IsAdmin"]; ok {
			user.IsAdmin = isAdmin.(bool)
		}
	}

	return user, nil
}
