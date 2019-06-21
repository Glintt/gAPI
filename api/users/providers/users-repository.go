package providers

import (
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/users/models"

)

const (
	PAGE_LENGTH      = 10
)

type UserRepository interface {
	InitUsers()
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
	FindUsersByUsernameOrEmail(q string, page int) []models.User
	GetUserByUsername(username string) []models.User
	GetUserByIdentifier(id string) models.User

	OpenTransaction() error
	CommitTransaction()
	RollbackTransaction()
	Release()
}

// NewUserRepository create an user repository based on the database
func NewUserRepository(user models.User) UserRepository {
	if database.SD_TYPE == "mongo" {
		session, db := database.GetSessionAndDB(database.MONGO_DB)
		collection := db.C(USERS_COLLECTION)

		return &UserMongoRepository{
			Session:    session,
			Db:         db,
			Collection: collection,
		}
	}
	if database.SD_TYPE == "oracle" {
		db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
		if err != nil {
			return nil
		}
		return &UserOracleRepository{
			Db:      db,
			DbError: err,
		}
	}
	return nil
}