package users

import (
	"database/sql"
	"gAPIManagement/api/database"
	"gAPIManagement/api/utils"

	"gopkg.in/mgo.v2/bson"
)

var INSERT_USER_ORACLE = `insert into gapi_users (id, username, password, email, isadmin) values (:id, :username, :password, :email, :isadmin)`
var UPDATE_USER_ORACLE = `UPDATE gapi_users
SET 
	username = :username, 
	password = :password,
	email = :email, 
	isadmin = :isadmin
WHERE id = :id
`
var FIND_BY_EMAIL_OR_USERNAME = `select id, username, password, email, isadmin from gapi_users where email like :email or username like :username`
var FIND_BY_USERNAME = `select id, username, password, email, isadmin from gapi_users where username = :username`

func InitUsersOracle() {
}

func CreateUserOracle(user User) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	user.Id = bson.NewObjectId()
	_, err = db.Exec(INSERT_USER_ORACLE,
		user.Id.Hex(), user.Username, user.Password, user.Email, utils.BoolToInt(user.IsAdmin),
	)

	database.CloseOracleConnection(db)
	if err != nil {
		return err
	}
	return err
}

func UpdateUserOracle(user User) error {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return err
	}

	_, err = db.Exec(UPDATE_USER_ORACLE,
		user.Username, user.Password, user.Email, utils.BoolToInt(user.IsAdmin), user.Id.Hex(),
	)

	database.CloseOracleConnection(db)
	if err != nil {
		return err
	}
	return err
}

func FindUsersByUsernameOrEmailOracle(q string, page int) []User {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []User{}
	}

	rows, err := db.Query(FIND_BY_EMAIL_OR_USERNAME, "%"+q+"%", "%"+q+"%")
	if err != nil {
		database.CloseOracleConnection(db)
		return []User{}
	}

	users := RowsToUser(rows)

	database.CloseOracleConnection(db)
	return users
}

func GetUserByUsernameOracle(username string) []User {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return []User{}
	}

	rows, err := db.Query(FIND_BY_USERNAME, username)
	if err != nil {
		database.CloseOracleConnection(db)
		return []User{}
	}

	users := RowsToUser(rows)

	database.CloseOracleConnection(db)
	return users
}

func RowsToUser(rows *sql.Rows) []User {
	var users []User
	for rows.Next() {
		var user User
		var id string
		rows.Scan(&id, &user.Username, &user.Password, &user.Email, &user.IsAdmin)

		if bson.IsObjectIdHex(id) {
			user.Id = bson.ObjectIdHex(id)
		} else {
			user.Id = bson.NewObjectId()
		}
		users = append(users, user)
	}

	defer rows.Close()

	return users
}
