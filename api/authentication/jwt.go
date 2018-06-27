package authentication

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gAPIManagement/api/users"
	"github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
	"gAPIManagement/api/config"
	"errors"
	"strings"
	"time"
)

var MinSizeSigningKey = 10
var MinExpirationTime = 30

var SIGNING_KEY = "AllYourBase"
var EXPIRATION_TIME = MinExpirationTime
var SERVICE_NAME = "authentication"

type TokenRequestObj struct{
	Username string `json:username`
	Password string `json:password`
}

type TokenCustomClaims struct {
	Username string `json:"Username"`
    jwt.StandardClaims
}


func InitGAPIAuthenticationServer(){
	if config.GApiConfiguration.Authentication.TokenExpirationTime > MinExpirationTime {
		EXPIRATION_TIME = config.GApiConfiguration.Authentication.TokenExpirationTime 
	}
	if len(config.GApiConfiguration.Authentication.TokenSigningKey) > MinSizeSigningKey {
		SIGNING_KEY = config.GApiConfiguration.Authentication.TokenSigningKey
	}
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

	if config.GApiConfiguration.Authentication.LDAP.Active && AuthenticateWithLDAP(username, password) {
		user := users.User{
			Username: strings.Split(username, "@")[0],
			Password: password,
			Email: username,
		}
		users.CreateUser(user)
		
		return user, nil
	}

	userList := users.GetUserByUsername(username)
	if len(userList) == 0 {
		return users.User{}, errors.New("Not Authorized.")
	}

	user := userList[0]
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		fmt.Println(err.Error())
	}
	if username == user.Username && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		
		return user, nil
	}

	return users.User{}, errors.New("Not Authorized.")
}


func NotAuthorized(c *routing.Context) error {
	c.Response.SetBody([]byte(`{"error":true, "msg": "Not authorized."}`))
	c.Response.SetStatusCode(401)
	c.Response.Header.SetContentType("application/json")
	return errors.New("Not allowed")
}

func AuthorizationMiddleware(c *routing.Context) error {
	token := c.Request.Header.Peek("Authorization")

	userClaims, validate := ValidateToken(string(token))

	if validate != nil {
		NotAuthorized(c)
		c.Abort()
		return nil
	}

	c.Request.Header.Add("User", userClaims["Username"].(string))
	return nil
}

func UserNotAllowed(c *routing.Context) error {
	c.Response.SetBody([]byte(`{"error":true, "msg": "Not authorized to access resource."}`))
	c.Response.SetStatusCode(405)
	c.Response.Header.SetContentType("application/json")
	return errors.New("Not allowed")
}

func AdminRequiredMiddleware(c *routing.Context) error {
	token := c.Request.Header.Peek("Authorization")

	claims, validate := ValidateToken(string(token))

	if validate != nil {
		NotAuthorized(c)
		c.Abort()
		return nil
	}

	username := claims["Username"].(string)
	usersList := users.GetUserByUsername(username)

	if len(usersList) == 0 || len(usersList) > 1 || !usersList[0].IsAdmin {
		UserNotAllowed(c)
		c.Abort()
		return nil
	}
	
	c.Request.Header.Add("User", username)
	return nil
}