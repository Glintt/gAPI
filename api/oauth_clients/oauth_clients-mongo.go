package oauth_clients

import (
	"fmt"

	"github.com/Glintt/gAPI/api/database"

	"gopkg.in/mgo.v2/bson"
)

func FindOAuthClientMongo(clientId string, clientSecret string) OAuthClient {
	session, db := database.GetSessionAndDB(database.MONGO_DB)
	fmt.Println(database.MONGO_DB)
	query := bson.M{"clientid": clientId, "clientsecret": clientSecret}

	oauthClient := OAuthClient{}

	err := db.C(OAUTH_CLIENTS_COLLECTION).Find(query).One(&oauthClient)
	if err != nil {
		return OAuthClient{}
	}

	database.MongoDBPool.Close(session)

	return oauthClient
}
