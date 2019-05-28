package routes

import (
	"github.com/Glintt/gAPI/api/users"
	"github.com/Glintt/gAPI/api/controllers"
	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/config"
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
	usersGroup.Put("/<username>", authentication.AuthorizationMiddleware, controllers.UpdateUserHandler)
	usersGroup.Put("/admin/<username>", authentication.AdminRequiredMiddleware, controllers.UpdateUserByAdminHandler)
}