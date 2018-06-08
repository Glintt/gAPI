package controllers

import (
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
		http.Response(c, `{"error" : false, "msg": "User created successfuly."}`, 200, UsersServiceName())		
	}
	
	return nil
}
