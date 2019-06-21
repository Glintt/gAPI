package providers

import (
	"database/sql"

	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/oauth_clients/models"
)

const (
	FIND_CLIENT = `select client_id, client_secret, user_id from gapi_oauth_clients where client_id = :client_id and client_secret = :client_secret`
	FIND_CLIENT_BY_CLIENT_ID = `select client_id, client_secret, user_id from gapi_oauth_clients where client_id = :client_id`
	FIND_FOR_USER = `select client_id, client_secret, user_id from gapi_oauth_clients where user_id = :user_id`
)


type OAuthClientOracleRepository struct {
	Db      *sql.DB
	DbError error
	Tx      *sql.Tx
}

// OpenTransaction open new database transaction
func (oauthClientRepos *OAuthClientOracleRepository) OpenTransaction() error {
	tx, err := oauthClientRepos.Db.Begin()
	oauthClientRepos.Tx = tx
	return err
}

// CommitTransaction commit database transaction
func (oauthClientRepos *OAuthClientOracleRepository) CommitTransaction() {
	oauthClientRepos.Tx.Commit()
}

// RollbackTransaction rollback database transaction
func (oauthClientRepos *OAuthClientOracleRepository) RollbackTransaction() {
	oauthClientRepos.Tx.Rollback()
}

func (oauthClientRepos *OAuthClientOracleRepository) Release() {
	database.CloseOracleConnection(oauthClientRepos.Db)	
}

func (oauthClientRepos *OAuthClientOracleRepository) FindOAuthClient(clientId string, clientSecret string) models.OAuthClient {
	rows, err := oauthClientRepos.Db.Query(FIND_CLIENT, clientId, clientSecret)

	if err != nil {
		return models.OAuthClient{}
	}

	clients := RowsToClients(rows, false)

	if len(clients) == 0 {
		return models.OAuthClient{}
	}
	return clients[0]
}

func (oauthClientRepos *OAuthClientOracleRepository) FindOAuthClientForUser(userID string) []models.OAuthClient {
	rows, err := oauthClientRepos.Db.Query(FIND_FOR_USER, userID)

	if err != nil {
		return []models.OAuthClient{}
	}

	clients := RowsToClients(rows, false)
	
	if len(clients) == 0 {
		return []models.OAuthClient{}
	}
	return clients
}

func (oauthClientRepos *OAuthClientOracleRepository) FindOAuthClientByClientId(clientID string) []models.OAuthClient {
	rows, err := oauthClientRepos.Db.Query(FIND_CLIENT_BY_CLIENT_ID, clientID)

	if err != nil {
		return []models.OAuthClient{}
	}

	clients := RowsToClients(rows, false)
	
	if len(clients) == 0 {
		return []models.OAuthClient{}
	}
	return clients
}

func RowsToClients(rows *sql.Rows, containsPagination bool) []models.OAuthClient {
	var clients []models.OAuthClient
	for rows.Next() {
		var client models.OAuthClient
		var r int
		if containsPagination {
			rows.Scan(&client.ClientId, &client.ClientSecret, &client.UserId, &r)
		} else {
			rows.Scan(&client.ClientId, &client.ClientSecret, &client.UserId)
		}

		clients = append(clients, client)
	}

	defer rows.Close()

	return clients
}
