package routes

import (
	"github.com/Glintt/gAPI/api/controllers"
	"github.com/Glintt/gAPI/api/authentication"
	"github.com/qiangxue/fasthttp-routing"
)

func InitGApiRoutes(router *routing.Router) {
	
	router.Get("/reload", authentication.AdminRequiredMiddleware, controllers.ReloadServices)
	router.Get("/invalidate-cache", authentication.AdminRequiredMiddleware, controllers.InvalidateCache)
	router.Get("/profile-gapi", controllers.ProfileGApiUsage)
}