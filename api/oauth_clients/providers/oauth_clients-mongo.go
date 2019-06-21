package providers

import (
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/oauth_clients/models"

	
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OAuthClientMongoRepository struct {
	Session    *mgo.Session
	Db         *mgo.Database
	Collection *mgo.Collection
}

// OpenTransaction open new database transaction
func (oauthClientRepos *OAuthClientMongoRepository) OpenTransaction() error {
	return nil
}

// CommitTransaction commit database transaction
func (oauthClientRepos *OAuthClientMongoRepository) CommitTransaction() {}

// RollbackTransaction rollback database transaction
func (oauthClientRepos *OAuthClientMongoRepository) RollbackTransaction() {}

// Release releases database connection to the pool
func (oauthClientRepos *OAuthClientMongoRepository) Release() {
	database.MongoDBPool.Close(oauthClientRepos.Session)
}

func (oauthClientRepos *OAuthClientMongoRepository) FindOAuthClient(clientId string, clientSecret string) models.OAuthClient {
	query := bson.M{"clientid": clientId, "clientsecret": clientSecret}

	oauthClient := models.OAuthClient{}

	err := oauthClientRepos.Collection.Find(query).One(&oauthClient)
	if err != nil {
		return models.OAuthClient{}
	}

	return oauthClient
}

func (oauthClientRepos *OAuthClientMongoRepository) FindOAuthClientForUser(userId string) []models.OAuthClient {
	query := bson.M{"userid": userId}

	oauthClient := []models.OAuthClient{}

	err := oauthClientRepos.Collection.Find(query).All(&oauthClient)
	if err != nil {
		return []models.OAuthClient{}
	}

	return oauthClient
}

func (oauthClientRepos *OAuthClientMongoRepository) FindOAuthClientByClientId(clientID string) []models.OAuthClient {
	query := bson.M{"clientid": clientID}

	oauthClient := []models.OAuthClient{}

	err := oauthClientRepos.Collection.Find(query).All(&oauthClient)
	if err != nil {
		return []models.OAuthClient{}
	}

	return oauthClient
}