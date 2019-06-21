package providers

import (
	"github.com/Glintt/gAPI/api/oauth_clients/models"
	"github.com/Glintt/gAPI/api/database"
	userModels "github.com/Glintt/gAPI/api/users/models"

)
const (
	OAUTH_CLIENTS_COLLECTION = "gapi_oauth_clients"
	PAGE_LENGTH              = 10
)

type OAuthClientRepository interface {
	FindOAuthClient(clientId string, clientSecret string) models.OAuthClient
	FindOAuthClientByClientId(clientId string) []models.OAuthClient
	FindOAuthClientForUser(userId string) []models.OAuthClient

	
	OpenTransaction() error
	CommitTransaction()
	RollbackTransaction()
	Release()
}



// NewServiceGroupRepository create an application group repository based on the database
func NewOAuthClientRepository(user userModels.User) OAuthClientRepository {
	if database.SD_TYPE == "mongo" {
		session, db := database.GetSessionAndDB(database.MONGO_DB)
		collection := db.C(OAUTH_CLIENTS_COLLECTION)

		return &OAuthClientMongoRepository{
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
		return &OAuthClientOracleRepository{
			Db:      db,
			DbError: err,
		}
	}
	return nil
}