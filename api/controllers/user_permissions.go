package controllers

import (
	"encoding/json"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/users"
	"github.com/Glintt/gAPI/api/user_permission"
	user_permission_models "github.com/Glintt/gAPI/api/user_permission/models"
	
	routing "github.com/qiangxue/fasthttp-routing"
)

func PermissionsServiceName() string {
	return user_permission.SERVICE_NAME
}


func GetUserPermissionsHandler(c *routing.Context) error {
	user := users.GetUserByUsername(c.Param("username"))
	
	if len(user) == 0 {
		http.Response(c, `{"error" : true, "msg": "User not found."}`, 404, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}

	userPermissions, err := user_permission.GetUserPermissions(user[0].Id.Hex())
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "Error getting user permissions: `+ err.Error()+`"}`, 500, PermissionsServiceName(), config.APPLICATION_JSON)
		return nil
	}
	userPermissionsJson, _ := json.Marshal(userPermissions)

	http.Response(c, string(userPermissionsJson), 200, PermissionsServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateUserPermissionHandler(c *routing.Context) error {
	user := users.GetUserByUsername(c.Param("username"))
	
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
