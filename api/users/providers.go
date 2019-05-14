package users

var UserMethods = map[string]map[string]interface{}{
	"mongo": {
		"init":              InitUsersMongo,
		"create":            CreateUserMongo,
		"update":            UpdateUserMongo,
		"findbyuseroremail": FindUsersByUsernameOrEmailMongo,
		"findbyusername":    GetUserByUsernameMongo},
	"oracle": {
		"init":              InitUsersOracle,
		"create":            CreateUserOracle,
		"update":            UpdateUserOracle,
		"findbyuseroremail": FindUsersByUsernameOrEmailOracle,
		"findbyusername":    GetUserByUsernameOracle}}
