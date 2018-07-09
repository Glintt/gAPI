package routes

import (
	"gAPIManagement/api/controllers"
	"gAPIManagement/api/authentication"
	"github.com/qiangxue/fasthttp-routing"
)

func InitGApiRoutes(router *routing.Router) {
	
	router.Get("/reload", authentication.AdminRequiredMiddleware, controllers.ReloadServices)
	router.Get("/invalidate-cache", authentication.AdminRequiredMiddleware, controllers.InvalidateCache)
	router.Get("/profile-gapi", controllers.ProfileGApiUsage)
}