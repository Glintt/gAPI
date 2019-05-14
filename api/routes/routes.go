package routes

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

func InitAPIRoutes(router *routing.Router) {
	InitServiceDocumentationRoutes(router)
	InitGApiRoutes(router)
	InitUsersService(router)
	InitAuthenticationAPI(router)
	InitServiceDiscoveryAPIRoutes(router)
	InitAnalyticsAPI(router)
	InitAppsGroupsApi(router)
	InitPluginsRoutes(router)
}
