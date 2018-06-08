package routes

import (
	"gAPIManagement/api/users"
	"gAPIManagement/api/controllers"
	"gAPIManagement/api/authentication"
	"gAPIManagement/api/config"
	"github.com/qiangxue/fasthttp-routing"
)



func InitUsersService(router *routing.Router) {
	users.InitUsers()
	
	usersGroup := router.Group(config.USERS_GROUP)

	usersGroup.Get("", authentication.AdminRequiredMiddleware, controllers.FindUsersHandler)
	usersGroup.Get("/", authentication.AdminRequiredMiddleware, controllers.FindUsersHandler)
	usersGroup.Post("/", authentication.AdminRequiredMiddleware, controllers.CreateUserHandler)
	usersGroup.Post("", authentication.AdminRequiredMiddleware, controllers.CreateUserHandler)
	usersGroup.Get("/<username>", authentication.AdminRequiredMiddleware, controllers.GetUserHandler)
}