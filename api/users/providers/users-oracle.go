package providers

import (
	"database/sql"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/utils"
	"github.com/Glintt/gAPI/api/users/models"

	"gopkg.in/mgo.v2/bson"
)

type UserOracleRepository struct {
	Db      *sql.DB
	DbError error
	Tx      *sql.Tx
}

const (
	FIND_BY_EMAIL_OR_USERNAME = `select id, username, password, email, isadmin from gapi_users where email like :email or username like :username`
	FIND_BY_USERNAME = `select id, username, password, email, isadmin from gapi_users where username = :username`
	UPDATE_USER_ORACLE = `UPDATE gapi_users
	SET 
		username = :username, 
		password = :password,
		email = :email, 
		isadmin = :isadmin
	WHERE id = :id
	`
	INSERT_USER_ORACLE = `insert into gapi_users (id, username, password, email, isadmin) values (:id, :username, :password, :email, :isadmin)`
)

// InitUsers inits user required modules
func (agmr *UserOracleRepository) InitUsers() {}

// OpenTransaction opens a new database transaction
func (agmr *UserOracleRepository) OpenTransaction() error {
	tx, err := agmr.Db.Begin()
	agmr.Tx = tx
	return err
}
// CommitTransaction commits a database transaction
func (agmr *UserOracleRepository) CommitTransaction() {
	agmr.Tx.Commit()
}
// RollbackTransaction rollbacks a database transaction
func (agmr *UserOracleRepository) RollbackTransaction() {
	agmr.Tx.Rollback()
}
// Release releases a database connection to the pool
func (agmr *UserOracleRepository) Release() {
	database.CloseOracleConnection(agmr.Db)
}

// CreateUser create a new user
func (ur *UserOracleRepository) CreateUser(user models.User) error {
	user.Id = bson.NewObjectId()
	_, err := ur.Tx.Exec(INSERT_USER_ORACLE,
		user.Id.Hex(), user.Username, user.Password, user.Email, utils.BoolToInt(user.IsAdmin),
	)
	return err
}

// UpdateUser update an existing user
func (ur *UserOracleRepository) UpdateUser(user models.User) error {
	_, err := ur.Tx.Exec(UPDATE_USER_ORACLE,
		user.Username, user.Password, user.Email, utils.BoolToInt(user.IsAdmin), user.Id.Hex(),
	)
	return err
}

// FindUsersByUsernameOrEmail return a list of users searched by username or email
func (ur *UserOracleRepository) FindUsersByUsernameOrEmail(q string, page int) []models.User {
	query := `SELECT * FROM
		(
			SELECT a.*, rownum r__
			FROM
			(
				` + FIND_BY_EMAIL_OR_USERNAME + `
			) a
			WHERE rownum < ((:page * 10) + 1 )
		)
		WHERE r__ >= (((:page-1) * 10) + 1)`

	rows, err := ur.Tx.Query(query, "%"+q+"%", "%"+q+"%", page)

	if err != nil {
		return []models.User{}
	}

	return RowsToUser(rows, true)
}

// GetUserByUsername return a list of users search by username
func (ur *UserOracleRepository) GetUserByUsername(username string) []models.User {
	rows, err := ur.Tx.Query(FIND_BY_USERNAME, username)
	if err != nil {
		return []models.User{}
	}

	return RowsToUser(rows, false)
}

func RowsToUser(rows *sql.Rows, containsPagination bool) []models.User {
	var users []models.User
	for rows.Next() {
		var user models.User
		var id string
		var r int
		if containsPagination {
			rows.Scan(&id, &user.Username, &user.Password, &user.Email, &user.IsAdmin, &r)
		} else {
			rows.Scan(&id, &user.Username, &user.Password, &user.Email, &user.IsAdmin)
		}

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
