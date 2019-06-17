package routes

import (
	"github.com/Glintt/gAPI/api/controllers"
	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/config"
	"github.com/qiangxue/fasthttp-routing"
)



func InitUserPermissionsService(router *routing.Router) {
	userPermissionsGroup := router.Group(config.USER_PERMISSIONS_GROUP)

	userPermissionsGroup.Get("/<username>", authentication.AdminRequiredMiddleware, controllers.GetUserPermissionsHandler)
	userPermissionsGroup.Put("/<username>", authentication.AdminRequiredMiddleware, controllers.UpdateUserPermissionHandler)
	userPermissionsGroup.Post("/<username>/<application_id>", authentication.AdminRequiredMiddleware, controllers.AddPermissionToApplicationGroupHandler)	
	userPermissionsGroup.Delete("/<username>/<application_id>", authentication.AdminRequiredMiddleware, controllers.RemovePermissionFromApplicationGroupHandler)	
}