package oauth_clients

var OAuthClientsMethods = map[string]map[string]interface{}{
	"mongo": {
		"find": FindOAuthClientMongo,
	},
	"oracle": {
		"find": FindOAuthClientOracle}}
