package authentication

import (
	"errors"
	"strings"
	"time"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/oauth_clients"
	"github.com/Glintt/gAPI/api/users"
	userModels "github.com/Glintt/gAPI/api/users/models"
	"github.com/Glintt/gAPI/api/utils"

	jwt "github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
	"golang.org/x/crypto/bcrypt"
)

var MinSizeSigningKey = 10
var MinExpirationTime = 30

var SIGNING_KEY = "AllYourBase"
var EXPIRATION_TIME = MinExpirationTime
var SERVICE_NAME = "authentication"

type TokenRequestObj struct {
	Username string `json:username`
	Password string `json:password`
}

type TokenCustomClaims struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}

func getUserService() users.UserService {
	uServ, err := users.NewUserServiceWithUser(userModels.GetInternalAPIUser())
	if err != nil {
		utils.LogMessage("Error creating UserService object", utils.InfoLogType)
	}
	return uServ
}

func InitGAPIAuthenticationServer() {
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

func GenerateToken(username string, password string) (string, error) {
	user, err := ValidateUserCredentials(username, password)
	if err != nil {
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

func LdapUserCreateOrUpdate(user userModels.User) {
	userService := getUserService()
	err := userService.CreateUser(user)
	if err != nil {
		utils.LogMessage("Create user error = "+err.Error(), utils.DebugLogType)

		userList := userService.GetUserByUsername(user.Username)
		if len(userList) == 0 {
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(userList[0].Password), []byte(user.Password)) != nil {
			user.Id = userList[0].Id
			user.IsAdmin = userList[0].IsAdmin

			hashedPwd, _ := userModels.GeneratePassword(user.Password)
			user.Password = string(hashedPwd)

			err = userService.UpdateUser(user)
			if err != nil {
				utils.LogMessage("Update user error = "+err.Error(), utils.DebugLogType)
			}
		}
	}
}

func ValidateUserCredentials(username string, password string) (userModels.User, error) {

	if config.GApiConfiguration.Authentication.LDAP.Active && AuthenticateWithLDAP(username, password) {
		email := username
		username = strings.Split(username, "@")[0]

		user := userModels.User{
			Username: username,
			Password: password,
			Email:    email,
		}

		LdapUserCreateOrUpdate(user)
	}

	userService := getUserService()
	
	userList := userService.GetUserByUsername(username)
	if len(userList) == 0 {
		return userModels.User{}, errors.New("Not Authorized.")
	}

	user := userList[0]
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		utils.LogMessage("Compare Password Hash Error = "+err.Error(), utils.DebugLogType)
	}
	if username == user.Username && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		return user, nil
	}

	return userModels.User{}, errors.New("Not Authorized.")
}

func NotAuthorized(c *routing.Context) error {
	c.Response.SetBody([]byte(`{"error":true, "msg": "Not authorized."}`))
	c.Response.SetStatusCode(401)
	c.Response.Header.SetContentType("application/json")
	return errors.New("Not allowed")
}

func GetUserFromToken(c *routing.Context) (jwt.MapClaims, error) {
	token := c.Request.Header.Peek("Authorization")

	return ValidateToken(string(token))
}

func CheckUserMiddleware(c *routing.Context) error {
	userClaims, validate := GetUserFromToken(c)

	if validate == nil {
		username := userClaims["Username"].(string)
		userService := getUserService()
	
		user := userService.GetUserByUsername(username)
		if len(user) == 0 {
			NotAuthorized(c)
			c.Abort()
			return nil
		}
		//userJson, _ := json.Marshal(user[0])
		c.Set("User", user[0])
		//c.Request.Header.Add("User", string(userJson))
	}

	return nil
}

func AuthorizationMiddleware(c *routing.Context) error {
	userClaims, validate := GetUserFromToken(c)

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
	userService := getUserService()
	usersList := userService.GetUserByUsername(username)

	if len(usersList) == 0 || len(usersList) > 1 || !usersList[0].IsAdmin {
		UserNotAllowed(c)
		c.Abort()
		return nil
	}

	c.Request.Header.Add("User", username)
	return nil
}

func OAuthClientRequiredMiddleware(c *routing.Context) error {
	clientId := string(c.Request.Header.Peek("ClientId"))
	clientSecret := string(c.Request.Header.Peek("ClientSecret"))

	oauthClient := oauth_clients.Find(clientId, clientSecret)

	if oauthClient.ClientSecret != clientSecret {
		UserNotAllowed(c)
		c.Abort()
		return nil
	}

	c.Request.Header.Add("User", clientId+"_"+clientSecret)
	return nil
}


func GetAuthenticatedUser (c *routing.Context) userModels.User{
	userInt := c.Get("User")
	var user userModels.User
	if userInt != nil {
		user = userInt.(userModels.User)
	}

	return user;
}