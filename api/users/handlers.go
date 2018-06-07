package users

import (
	"strconv"
	"encoding/json"
	"gAPIManagement/api/http"
	"github.com/qiangxue/fasthttp-routing"
)

func GetUserHandler(c *routing.Context) error {
	user := GetUserByUsername(c.Param("username"))
	
	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, SERVICE_NAME)
		return nil
	}

	userJSON, _ := json.Marshal(user[0])
	http.Response(c, string(userJSON), 200, SERVICE_NAME)
	return nil
}


func FindUsersHandler(c *routing.Context) error {
	page := 1
	searchQuery := ""

	if c.QueryArgs().Has("page") {
		var err error
		page, err = strconv.Atoi(string(c.QueryArgs().Peek("page")))

		if err != nil {
			http.Response(c, `{"error" : true, "msg": "Invalid page provided."}`, 404, SERVICE_NAME)
			return nil
		}
	}
	if c.QueryArgs().Has("q") {
		searchQuery = string(c.QueryArgs().Peek("q"))
	}
	
	users := FindUsersByUsernameOrEmail(searchQuery, page)

	userJSON,_ := json.Marshal(users)

	if users == nil {
		userJSON = []byte(`[]`)
	}
	
	http.Response(c, string(userJSON), 200, SERVICE_NAME)
	return nil
}

func CreateUserHandler(c *routing.Context) error {
	var user User
	err := json.Unmarshal(c.Request.Body(), &user)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be created."}`, 400, SERVICE_NAME)
		return nil
	}

	err = CreateUser(user)
	
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "User could not be created."}`, 400, SERVICE_NAME)
	}else {
		http.Response(c, `{"error" : false, "msg": "User created successfuly."}`, 200, SERVICE_NAME)		
	}
	
	return nil
}
