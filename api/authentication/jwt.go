package authentication

import (
	"golang.org/x/crypto/bcrypt"
	"gAPIManagement/api/users"
	"github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
	"gAPIManagement/api/config"
	"strconv"
	"errors"
	"encoding/json"
	"strings"
	"time"
)

var MinSizeSigningKey = 10
var MinExpirationTime = 30

var SIGNING_KEY = "AllYourBase"
var EXPIRATION_TIME = MinExpirationTime

type TokenRequestObj struct{
	Username string `json:username`
	Password string `json:password`
}

type TokenCustomClaims struct {
	Username string `json:"Username"`
    jwt.StandardClaims
}


func InitGAPIAuthenticationServer(router *routing.Router){
	if config.GApiConfiguration.Authentication.TokenExpirationTime > MinExpirationTime {
		EXPIRATION_TIME = config.GApiConfiguration.Authentication.TokenExpirationTime 
	}
	if len(config.GApiConfiguration.Authentication.TokenSigningKey) > MinSizeSigningKey {
		SIGNING_KEY = config.GApiConfiguration.Authentication.TokenSigningKey
	}

	// LoadUsers()

	router.Post("/oauth/token", GetTokenHandler)
	router.Get("/oauth/authorize", AuthorizeTokenHandler)
	router.Get("/oauth/me", MeHandler)
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

func MeHandler(c *routing.Context) error {
	c.Response.Header.SetContentType("application/json")
	authorizationToken := c.Request.Header.Peek("Authorization")

	tokenClaims, err := ValidateToken(string(authorizationToken))

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

	_, err := ValidateToken(string(authorizationToken))

	if err != nil{
		c.Response.SetBody([]byte(`{"error":true, "msg":"`+ err.Error() + `"}`))
		c.Response.Header.SetStatusCode(401)
		return nil
	}
	c.Response.SetBody([]byte(`{"error":false, "msg":"Token is valid."}`))
	return nil
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	tokenParsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNING_KEY), nil
	})

	if err != nil {
		return nil, errors.New("Not Authorized.")
	}
	claims := tokenParsed.Claims.(jwt.MapClaims)
	
	if err == nil && tokenParsed.Valid {
		return claims, nil
	}

	return nil, errors.New("Not Authorized.")
}

func GenerateToken(username string, password string) (string, error){
	_, err := ValidateUserCredentials(username, password)
	if err != nil{
		return "", err
	}

	mySigningKey := []byte(SIGNING_KEY)

	// Create the Claims
	claims := TokenCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: (time.Now().Unix() + int64(EXPIRATION_TIME)),
			Issuer:    "gAPI",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	
	return ss, nil
}

func ValidateUserCredentials(username string, password string) (users.User, error) {
	userList := users.GetUserByUsername(username)
	if len(userList) == 0 {
		return users.User{}, errors.New("Not Authorized.")
	}

	user := userList[0]
	if username == user.Username && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		return user, nil
	}

	return users.User{}, errors.New("Not Authorized.")
}


func NotAllowed(c *routing.Context) error {
	c.Response.SetBody([]byte(`{"error":true, "msg": "Not authorized."}`))
	c.Response.SetStatusCode(401)
	c.Response.Header.SetContentType("application/json")
	return errors.New("Not allowed")
}

func AuthorizationMiddleware(c *routing.Context) error {
	token := c.Request.Header.Peek("Authorization")

	_, validate := ValidateToken(string(token))

	if validate != nil {
		NotAllowed(c)
		c.Abort()
		return nil
	}

	return nil
}

func AdminRequiredMiddleware(c *routing.Context) error {
	token := c.Request.Header.Peek("Authorization")

	claims, validate := ValidateToken(string(token))

	if validate != nil {
		NotAllowed(c)
		c.Abort()
		return nil
	}

	username := claims["Username"].(string)
	usersList := users.GetUserByUsername(username)

	if len(usersList) == 0 || len(usersList) > 1 || !usersList[0].IsAdmin {
		NotAllowed(c)
		c.Abort()
		return nil
	}

	return nil
}