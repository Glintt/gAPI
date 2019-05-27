package routes

import (
	"gAPIManagement/api/authentication"
	"gAPIManagement/api/config"
	"gAPIManagement/api/controllers"

	routing "github.com/qiangxue/fasthttp-routing"
)

func InitAnalyticsAPI(router *routing.Router) {
	analyticsAPI := router.Group(config.ANALYTICS_GROUP)

	analyticsAPI.Get("/applications", authentication.AdminRequiredMiddleware, controllers.ApplicationGroupAnalytics)
	analyticsAPI.Get("/api", authentication.AdminRequiredMiddleware, controllers.APIAnalytics)
	analyticsAPI.Get("/logs", authentication.AdminRequiredMiddleware, controllers.Logs)
}
