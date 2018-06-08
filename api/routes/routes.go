package routes

import (
	"github.com/qiangxue/fasthttp-routing"
)

func InitAPIRoutes(router *routing.Router) {
	InitGApiRoutes(router)
	InitUsersService(router)
	InitAuthenticationAPI(router)
	InitServiceDiscoveryAPIRoutes(router)
	InitAnalyticsAPI(router)
}