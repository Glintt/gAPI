package routes

import (
	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/controllers"

	"github.com/qiangxue/fasthttp-routing"
)

func InitPluginsRoutes(router *routing.Router) {
	router.Get("/plugins", authentication.AdminRequiredMiddleware, controllers.ListPluginsAvailable)
	router.Get("/plugins/active", authentication.AdminRequiredMiddleware, controllers.ActivePlugins)
}
