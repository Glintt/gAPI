package controllers

import (
	"gAPIManagement/api/utils"
	"gAPIManagement/api/users"
	"strconv"
	"encoding/json"
	"gAPIManagement/api/http"
	"github.com/qiangxue/fasthttp-routing"
)
func UsersServiceName() string {
	return users.SERVICE_NAME
}
func GetUserHandler(c *routing.Context) error {
	user := users.GetUserByUsername(c.Param("username"))
	
	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, UsersServiceName())
		return nil
	}

	user[0].Password = ""
	userJSON, _ := json.Marshal(user[0])
	http.Response(c, string(userJSON), 200, UsersServiceName())
	return nil
}


func FindUsersHandler(c *routing.Context) error {
	page := 1
	searchQuery := ""

	if c.QueryArgs().Has("page") {
		var err error
		page, err = strconv.Atoi(string(c.QueryArgs().Peek("page")))

		if err != nil {
			http.Response(c, `{"error" : true, "msg": "Invalid page provided."}`, 404, UsersServiceName())
			return nil
		}
	}
	if c.QueryArgs().Has("q") {
		searchQuery = string(c.QueryArgs().Peek("q"))
	}
	
	users := users.FindUsersByUsernameOrEmail(searchQuery, page)

	userJSON,_ := json.Marshal(users)

	if users == nil {
		userJSON = []byte(`[]`)
	}
	
	http.Response(c, string(userJSON), 200, UsersServiceName())
	return nil
}

func CreateUserHandler(c *routing.Context) error {
	var user users.User
	err := json.Unmarshal(c.Request.Body(), &user)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be created."}`, 400, UsersServiceName())
		return nil
	}

	err = users.CreateUser(user)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be created."}`, 400, UsersServiceName())
	}else {
		http.Response(c, `{"error" : false, "msg": "User created successfuly."}`, 201, UsersServiceName())		
	}
	
	return nil
}

func UpdateUserByAdminHandler(c *routing.Context) error {
	usersList := users.GetUserByUsername(c.Param("username"))	
	
	if len(usersList) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, UsersServiceName())
		return nil
	}
	
	user := usersList[0]

	allowedToUpdate := []string{"Email", "Password", "IsAdmin"}
	user, err := updateUserUsingUpdateBody(user, c.Request.Body(), allowedToUpdate)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be updated."}`, 400, UsersServiceName())
		return nil
	}

	err = users.UpdateUser(user)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be updated."}`, 400, UsersServiceName())
	}else {
		http.Response(c, `{"error" : false, "msg": "User updated successfuly."}`, 200, UsersServiceName())		
	}
	return nil
}

func UpdateUserHandler(c *routing.Context) error {
	requestUser := string(c.Request.Header.Peek("User"))
	
	if requestUser != c.Param("username") {
		http.Response(c, `{"error" : true, "msg": "User could not be updated."}`, 400, UsersServiceName())
		return nil
	}

	usersList := users.GetUserByUsername(c.Param("username"))	
	
	if len(usersList) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, UsersServiceName())
		return nil
	}
	
	user := usersList[0]

	allowedToUpdate := []string{"email", "password"}
	user, err := updateUserUsingUpdateBody(user, c.Request.Body(), allowedToUpdate)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be updated."}`, 400, UsersServiceName())
		return nil
	}

	err = users.UpdateUser(user)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be updated."}`, 400, UsersServiceName())
	}else {
		http.Response(c, `{"error" : false, "msg": "User updated successfuly."}`, 200, UsersServiceName())		
	} 
	return nil
}


func updateUserUsingUpdateBody(user users.User, body []byte, allowedToUpdate []string) (users.User, error) {
	var userUpdateObj map[string]interface{}
	err := json.Unmarshal(body, &userUpdateObj)

	if err != nil {
		return user, err
	}

	if email, ok := userUpdateObj["Email"]; ok {
		user.Email = email.(string)
	}
	if pwd, ok := userUpdateObj["Password"]; ok {
		hashedPwd, _ := users.GeneratePassword(pwd.(string))
		user.Password = string(hashedPwd)
	}

	if utils.ArrayContainsString(allowedToUpdate, "IsAdmin") {
		if isAdmin, ok := userUpdateObj["IsAdmin"]; ok {
			user.IsAdmin = isAdmin.(bool)
		}
	}

	return user, nil
}