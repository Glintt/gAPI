package controllers

import (
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
		c.Response.SetBody([]byte(`{"error":true, "msg":"` + err.Error() + `"}`))
		c.Response.SetStatusCode(401)
		return nil
	}

	c.Response.SetBody([]byte(`{"token":"` + token + `", "expiration_time": ` + strconv.Itoa(auth.EXPIRATION_TIME) +`}`))
	return nil
}

func MeHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	authorizationToken := c.Request.Header.Peek("Authorization")

	tokenClaims, err := auth.ValidateToken(string(authorizationToken))

	if err != nil{
		c.Response.SetBody([]byte(`{"error":true, "msg":"`+ err.Error() + `"}`))
		c.Response.Header.SetStatusCode(400)
		return nil
	}

	username := tokenClaims["Username"].(string)
	usersList := users.GetUserByUsername(username)

	if len(usersList) == 0 || len(usersList) > 1 || !usersList[0].IsAdmin {
		c.Response.SetBody([]byte(`{"error":true, "msg":"`+ err.Error() + `"}`))
		c.Response.Header.SetStatusCode(400)
		return nil
	}

	userJSON,_ := json.Marshal(usersList[0])
	c.Response.SetBody(userJSON)
	return nil
}

func AuthorizeTokenHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	authorizationToken := c.Request.Header.Peek("Authorization")

	_, err := auth.ValidateToken(string(authorizationToken))

	if err != nil{
		c.Response.SetBody([]byte(`{"error":true, "msg":"`+ err.Error() + `"}`))
		c.Response.Header.SetStatusCode(401)
		return nil
	}
	c.Response.SetBody([]byte(`{"error":false, "msg":"Token is valid."}`))
	return nil
}