package routes

import (
	"gAPIManagement/api/cache"
	"gAPIManagement/api/controllers"
	"gAPIManagement/api/authentication"
	"gAPIManagement/api/config"
	"github.com/qiangxue/fasthttp-routing"
)


func InitServiceDiscoveryAPIRoutes(router *routing.Router) {
	serviceDiscoveryAPIGroup := router.Group(config.SERVICE_DISCOVERY_GROUP)

	serviceDiscoveryAPIGroup.Post("/register", authentication.AdminRequiredMiddleware, controllers.RegisterHandler)
	serviceDiscoveryAPIGroup.Post("/admin/normalize", authentication.AdminRequiredMiddleware, controllers.NormalizeServices)
	serviceDiscoveryAPIGroup.Post("/update", authentication.AdminRequiredMiddleware, controllers.UpdateHandler)
	serviceDiscoveryAPIGroup.Get("/services", cache.ResponseCacheGApi, controllers.ListServicesHandler, cache.StoreCacheGApi)
	serviceDiscoveryAPIGroup.Get("/endpoint", cache.ResponseCacheGApi, controllers.GetEndpointHandler, cache.StoreCacheGApi)
	serviceDiscoveryAPIGroup.Delete("/delete", authentication.AdminRequiredMiddleware, controllers.DeleteEndpointHandler)
	serviceDiscoveryAPIGroup.Post("/services/manage", authentication.AuthorizationMiddleware, controllers.ManageServiceHandler)
	serviceDiscoveryAPIGroup.Get("/services/manage/types", controllers.ManageServiceTypesHandler)
	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" {
		LoadDBSpecificEndpoints(serviceDiscoveryAPIGroup)
	}
	serviceDiscoveryAPIGroup.To("GET,POST,PUT,PATCH,DELETE", "/*", controllers.ServiceNotFound)
}

func LoadDBSpecificEndpoints(router *routing.RouteGroup) {
	router.Post("/service-groups", authentication.AdminRequiredMiddleware, controllers.RegisterServiceGroupHandler)
	router.Put("/service-groups/<group_id>", authentication.AdminRequiredMiddleware, controllers.UpdateServiceGroup)
	router.Delete("/service-groups/<group_id>", authentication.AdminRequiredMiddleware, controllers.RemoveServiceGroup)
	router.Post("/service-groups/<group_id>/services", authentication.AdminRequiredMiddleware, controllers.AddServiceToGroupHandler)
	router.Delete("/service-groups/<group_id>/services/<service_id>", authentication.AdminRequiredMiddleware, controllers.DeassociateServiceFromGroup)
	// sd.sdAPI.Post("/service-groups/service/register", authentication.AuthorizationMiddleware, RegisterServiceToServiceGroupHandler)
	router.Get("/service-groups", controllers.ListServiceGroupsHandler)
	router.Get("/service-groups/<group>", controllers.GetServiceGroupHandler)
}
