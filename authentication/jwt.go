package authentication

import (
	"github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
	"api-management/config"
	"strconv"
	"errors"
	"encoding/json"
	"strings"
	"time"
)
var SIGNING_KEY = "AllYourBase"
var EXPIRATION_TIME = 15000

type TokenRequestObj struct{
	Username string
	Password string
}

func InitGAPIAuthenticationServer(router *routing.Router){
	router.Post("/oauth/token", GetTokenHandler)
	router.Get("/oauth/authorize", AuthorizeTokenHandler)
}

func GetTokenHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	tokenRequestBody := c.Request.Body()
	var tokenRequestObj TokenRequestObj
	json.Unmarshal(tokenRequestBody, &tokenRequestObj)

	token, err := GenerateToken(tokenRequestObj.Username, tokenRequestObj.Password)

	if err != nil {
		c.Response.SetBody([]byte(`{"error":true, "msg":"` + err.Error() + `"}`))
		c.Response.SetStatusCode(401)
		return nil
	}

	c.Response.SetBody([]byte(`{"token":"` + token + `", "expiration_time": ` + strconv.Itoa(EXPIRATION_TIME) +`}`))
	return nil
}

func AuthorizeTokenHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	authorizationToken := c.Request.Header.Peek("Authorization")

	err := ValidateToken(string(authorizationToken))

	if err != nil{
		c.Response.SetBody([]byte(`{"error":true, "msg":"`+ err.Error() + `"}`))
		c.Response.Header.SetStatusCode(401)
		return nil
	}
	c.Response.SetBody([]byte(`{"error":false, "msg":"Token is valid."}`))
	return nil
}

func ValidateToken(tokenString string) error {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	tokenParsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNING_KEY), nil
	})

	if  err == nil && tokenParsed.Valid {
		return nil
	} 

	return errors.New("Not Authorized.")
}

func GenerateToken(username string, password string) (string, error){
	if err := ValidateUserCredentials(username, password); err != nil{
		return "", err
	}

	mySigningKey := []byte(SIGNING_KEY)

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: (time.Now().Unix() + int64(EXPIRATION_TIME)),
		Issuer:    "gAPI",
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	
	return ss, nil
}

func ValidateUserCredentials(username string, password string) error {
	if username == config.GApiConfiguration.Authentication.Username && password == config.GApiConfiguration.Authentication.Username {
		return nil
	}

	return errors.New("Not Authorized.")
}


func NotAllowed(c *routing.Context) error {
	c.Response.SetBody([]byte(`{"error":true, "msg": "Not authorized."}`))
	c.Response.SetStatusCode(401)
	c.Response.Header.SetContentType("application/json")
	return errors.New("Not allowed")
}

func AuthorizationMiddleware(c *routing.Context) error {
	token := c.Request.Header.Peek("Authorization")

	validate := ValidateToken(string(token))

	if validate != nil {
		NotAllowed(c)
		c.Abort()
		return nil
	}

	return nil
}