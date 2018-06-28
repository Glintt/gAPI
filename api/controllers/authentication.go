package controllers

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/authentication"
	auth "gAPIManagement/api/authentication"
	"gAPIManagement/api/users"
	routing "github.com/qiangxue/fasthttp-routing"
	"strconv"
	"encoding/json"
)

func GetTokenHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	tokenRequestBody := c.Request.Body()
	var tokenRequestObj auth.TokenRequestObj
	json.Unmarshal(tokenRequestBody, &tokenRequestObj)

	token, err := auth.GenerateToken(tokenRequestObj.Username, tokenRequestObj.Password)

	if err != nil {
		http.Response(c,`{"error":true, "msg":"` + err.Error() + `"}`, 401, authentication.SERVICE_NAME, config.APPLICATION_JSON)
		return nil
	}
	http.Response(c,`{"token":"` + token + `", "expiration_time": ` + strconv.Itoa(auth.EXPIRATION_TIME) +`}`, 200, authentication.SERVICE_NAME, config.APPLICATION_JSON)
	return nil
}

func MeHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	authorizationToken := c.Request.Header.Peek("Authorization")

	tokenClaims, err := auth.ValidateToken(string(authorizationToken))

	if err != nil{
		http.Response(c, `{"error":true, "msg":"`+ err.Error() + `"}`, 400, authentication.SERVICE_NAME, config.APPLICATION_JSON)
		return nil
	}

	username := tokenClaims["Username"].(string)
	usersList := users.GetUserByUsername(username)
	
	if len(usersList) == 0 || len(usersList) > 1 {
		http.Response(c, `{"error":true, "msg":"User not found."}`, 404, authentication.SERVICE_NAME, config.APPLICATION_JSON)
		return nil
	}

	usersList[0].Password = ""
	userJSON,_ := json.Marshal(usersList[0])
	http.Response(c, string(userJSON), 200, authentication.SERVICE_NAME, config.APPLICATION_JSON)
	return nil
}

func AuthorizeTokenHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	authorizationToken := c.Request.Header.Peek("Authorization")

	_, err := auth.ValidateToken(string(authorizationToken))

	if err != nil{
		http.Response(c, `{"error":true, "msg":"`+ err.Error() + `"}`, 401, authentication.SERVICE_NAME, config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error":false, "msg":"Token is valid."}`, 200, authentication.SERVICE_NAME, config.APPLICATION_JSON)

	return nil
}