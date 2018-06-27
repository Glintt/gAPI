package authentication

import (
	"gAPIManagement/api/utils"
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
	user, err := ValidateUserCredentials(username, password)
	if err != nil{
		return "", err
	}

	mySigningKey := []byte(SIGNING_KEY)

	// Create the Claims
	claims := TokenCustomClaims{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: (time.Now().Unix() + int64(EXPIRATION_TIME)),
			Issuer:    "gAPI",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	
	return ss, nil
}

func LdapUserCreateOrUpdate(user users.User) {
	err := users.CreateUser(user)
	if err != nil {
		utils.LogMessage("Create user error = " + err.Error(), utils.DebugLogType)

		userList := users.GetUserByUsername(user.Username)
		if len(userList) == 0 {
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(userList[0].Password), []byte(user.Password)) != nil {
			user.Id = userList[0].Id
			user.IsAdmin = userList[0].IsAdmin			

			hashedPwd, _ := users.GeneratePassword(user.Password)
			user.Password = string(hashedPwd)

			err = users.UpdateUser(user)
			if err != nil {
				utils.LogMessage("Update user error = " + err.Error(), utils.DebugLogType)
			}
		}
	}
}

func ValidateUserCredentials(username string, password string) (users.User, error) {

	if config.GApiConfiguration.Authentication.LDAP.Active && AuthenticateWithLDAP(username, password) {
		email := username
		username = strings.Split(username, "@")[0]

		user := users.User{
			Username: username,
			Password: password,
			Email: email,
		}
		
		LdapUserCreateOrUpdate(user)
	}
	
	userList := users.GetUserByUsername(username)
	if len(userList) == 0 {
		return users.User{}, errors.New("Not Authorized.")
	}

	user := userList[0]
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		utils.LogMessage("Compare Password Hash Error = " + err.Error(), utils.DebugLogType)
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