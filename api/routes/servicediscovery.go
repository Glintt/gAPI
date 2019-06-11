package routes

import (
	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/cache"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/controllers"

	routing "github.com/qiangxue/fasthttp-routing"
)

func InitServiceDiscoveryAPIRoutes(router *routing.Router) {
	serviceDiscoveryAPIGroup := router.Group(config.SERVICE_DISCOVERY_GROUP)

	serviceDiscoveryAPIGroup.Post("/services", authentication.AdminRequiredMiddleware, controllers.RegisterHandler)
	serviceDiscoveryAPIGroup.Post("/admin/normalize", authentication.AdminRequiredMiddleware, controllers.NormalizeServices)
	serviceDiscoveryAPIGroup.Put("/services/<service_id>", authentication.AdminRequiredMiddleware, controllers.UpdateHandler)
	serviceDiscoveryAPIGroup.Get("/services", authentication.CheckUserMiddleware, cache.ResponseCacheGApi, controllers.ListServicesHandler, cache.StoreCacheGApi)
	serviceDiscoveryAPIGroup.Get("/endpoint", cache.ResponseCacheGApi, authentication.CheckUserMiddleware, controllers.GetEndpointHandler, cache.StoreCacheGApi)
	serviceDiscoveryAPIGroup.Delete("/services/<service_id>", authentication.AdminRequiredMiddleware, controllers.DeleteEndpointHandler)
	serviceDiscoveryAPIGroup.Post("/services/manage", authentication.AuthorizationMiddleware, controllers.ManageServiceHandler)
	serviceDiscoveryAPIGroup.Get("/services/manage/types", controllers.ManageServiceTypesHandler)
	serviceDiscoveryAPIGroup.Post("/auto-register", authentication.OAuthClientRequiredMiddleware, controllers.AutoRegisterHandler)
	serviceDiscoveryAPIGroup.Post("/auto-deregister", authentication.OAuthClientRequiredMiddleware, controllers.AutoDeRegisterHandler)
	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" || config.GApiConfiguration.ServiceDiscovery.Type == "oracle" {
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
