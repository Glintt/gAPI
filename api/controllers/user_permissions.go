package controllers

import (
	"encoding/json"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery/appgroups"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	"github.com/Glintt/gAPI/api/user_permission"
	userModels "github.com/Glintt/gAPI/api/users/models"
	user_permission_models "github.com/Glintt/gAPI/api/user_permission/models"
	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

func PermissionsServiceName() string {
	return user_permission.SERVICE_NAME
}

func GetUserGroupsPermissionsHandler(c *routing.Context) error{
	userService := getUserService(c)
	user := userService.GetUserByUsername(c.Param("username"))

	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	appgroupService,_ := appgroups.NewApplicationGroupServiceWithUser(userModels.GetInternalAPIUser())
	applicationGroups, err := appgroupService.GetApplicationGroupsPermissions(user[0].Id.Hex())
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting user's permitted application groups: `+err.Error()+`"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}
	applicationGroupsJson, _ := json.Marshal(applicationGroups)
	http.Ok(c, string(applicationGroupsJson), PermissionsServiceName())
	return nil
}

func GetUserPermissionsHandler(c *routing.Context) error {
	userService := getUserService(c)
	user := userService.GetUserByUsername(c.Param("username"))

	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	userPermissions, err := user_permission.GetUserPermissions(user[0].Id.Hex())
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting user permissions: `+err.Error()+`"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}
	userPermissionsJson, _ := json.Marshal(userPermissions)

	http.Response(c, string(userPermissionsJson), 200, PermissionsServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateUserPermissionHandler(c *routing.Context) error {
	userService := getUserService(c)
	user := userService.GetUserByUsername(c.Param("username"))

	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}
	var newPermissions []user_permission_models.UserPermission

	err := json.Unmarshal(c.Request.Body(), &newPermissions)
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting permissions object"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	err = user_permission.UpdatePermission(user[0].Id.Hex(), newPermissions)
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error updating user permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "User permissions updated"}`, 201, PermissionsServiceName(), config.APPLICATION_JSON)
	return nil
}

func AddPermissionToApplicationGroupHandler(c *routing.Context) error {
	userService := getUserService(c)
	user := userService.GetUserByUsername(c.Param("username"))
	applicationId := c.Param("application_id")

	if !bson.IsObjectIdHex(applicationId) {
		http.Response(c, `{"error" : true, "msg": "Error getting user permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	// Get application group service
	appGroupService, err := getAppGroupService(c)
	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	userPermissions, err := user_permission.GetUserPermissions(user[0].Id.Hex())
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting user permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	group := appgroups.ApplicationGroup{Id: bson.ObjectIdHex(string(applicationId))}
	servicesList, err := appGroupService.GetServicesForApplicationGroup(group)
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting application group service's permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	for _, s := range servicesList {
		userPermissions = append(userPermissions, user_permission_models.UserPermission{
			UserId:    user[0].Id.Hex(),
			ServiceId: s.Id.Hex(),
		})
	}

	err = user_permission.UpdatePermission(user[0].Id.Hex(), userPermissions)
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error updating user permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "User permissions updated"}`, 201, PermissionsServiceName(), config.APPLICATION_JSON)
	return nil
}

func RemovePermissionFromApplicationGroupHandler(c *routing.Context) error {
	userService := getUserService(c)
	user := userService.GetUserByUsername(c.Param("username"))
	applicationId := c.Param("application_id")
	appGroupService, err := getAppGroupService(c)
	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	userPermissions, err := user_permission.GetUserPermissions(user[0].Id.Hex())
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting user permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	group := appgroups.ApplicationGroup{Id: bson.ObjectIdHex(string(applicationId))}
	servicesList, err := appGroupService.GetServicesForApplicationGroup(group)
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting application group service's permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	var finalUserPermissions []user_permission_models.UserPermission
	for _, u := range userPermissions {
		// If not on the list of group's services, add it to the user permissions
		if !ContainsService(servicesList, u) {
			finalUserPermissions = append(finalUserPermissions, u)
		}
	}

	err = user_permission.UpdatePermission(user[0].Id.Hex(), finalUserPermissions)
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error updating user permissions"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "User permissions updated"}`, 201, PermissionsServiceName(), config.APPLICATION_JSON)
	return nil
}

func ContainsService(servicesList []service.Service, permission user_permission_models.UserPermission) bool {
	for _, s := range servicesList {
		if s.Id.Hex() == permission.ServiceId {
			return true
		}
	}
	return false
}
