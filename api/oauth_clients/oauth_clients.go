package oauth_clients

import (
	"github.com/Glintt/gAPI/api/oauth_clients/providers"
	"github.com/Glintt/gAPI/api/oauth_clients/models"
	userModels "github.com/Glintt/gAPI/api/users/models"
)

const (
	SERVICE_NAME             = "gapi_oauth_clients"
)


func createRepository() providers.OAuthClientRepository{
	return providers.NewOAuthClientRepository(userModels.GetInternalAPIUser())
}

func Find(clientId string, clientSecret string) models.OAuthClient {
	repository := createRepository()
	return repository.FindOAuthClient(clientId, clientSecret)
}

func FindByClientId(clientId string) []models.OAuthClient {
	repository := createRepository()
	return repository.FindOAuthClientByClientId(clientId)
}

func FindForUser(userID string) []models.OAuthClient {
	repository := createRepository()
	return repository.FindOAuthClientForUser(userID)
}
