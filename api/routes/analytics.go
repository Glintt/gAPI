package routes

import (
	"gAPIManagement/api/controllers"
	"gAPIManagement/api/config"
	"gAPIManagement/api/authentication"
	"github.com/qiangxue/fasthttp-routing"
)

func InitAnalyticsAPI(router *routing.Router) {
	analyticsAPI := router.Group(config.ANALYTICS_GROUP)

	analyticsAPI.Get("/api", authentication.AdminRequiredMiddleware, controllers.APIAnalytics)
	analyticsAPI.Get("/logs", authentication.AdminRequiredMiddleware, controllers.Logs)
}