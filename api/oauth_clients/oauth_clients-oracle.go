package oauth_clients

import (
	"database/sql"

	"github.com/Glintt/gAPI/api/database"
)

var FIND_CLIENT = `select client_id, client_secret from gapi_oauth_clients where client_id = :client_id and client_secret = :client_secret`

func FindOAuthClientOracle(clientId string, clientSecret string) OAuthClient {
	db, err := database.ConnectToOracle(database.ORACLE_CONNECTION_STRING)
	if err != nil {
		return OAuthClient{}
	}

	rows, err := db.Query(FIND_CLIENT, clientId, clientSecret)

	if err != nil {
		database.CloseOracleConnection(db)
		return OAuthClient{}
	}

	clients := RowsToClients(rows, false)

	database.CloseOracleConnection(db)
	if len(clients) == 0 {
		return OAuthClient{}
	}
	return clients[0]
}

func RowsToClients(rows *sql.Rows, containsPagination bool) []OAuthClient {
	var clients []OAuthClient
	for rows.Next() {
		var client OAuthClient
		var r int
		if containsPagination {
			rows.Scan(&client.ClientId, &client.ClientSecret, &r)
		} else {
			rows.Scan(&client.ClientId, &client.ClientSecret)
		}

		clients = append(clients, client)
	}

	defer rows.Close()

	return clients
}
