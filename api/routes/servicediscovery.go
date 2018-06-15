package routes

import (
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
	serviceDiscoveryAPIGroup.Get("/services", controllers.ListServicesHandler)
	serviceDiscoveryAPIGroup.Get("/endpoint", controllers.GetEndpointHandler)
	serviceDiscoveryAPIGroup.Delete("/delete", authentication.AdminRequiredMiddleware, controllers.DeleteEndpointHandler)
	serviceDiscoveryAPIGroup.Post("/services/manage", authentication.AdminRequiredMiddleware, controllers.ManageServiceHandler)
	serviceDiscoveryAPIGroup.Get("/services/manage/types", controllers.ManageServiceTypesHandler)
	if config.GApiConfiguration.ServiceDiscovery.Type == "mongo" {
		LoadDBSpecificEndpoints(serviceDiscoveryAPIGroup)
	}
	serviceDiscoveryAPIGroup.To("GET,POST,PUT,PATCH,DELETE", "/*", controllers.ServiceNotFound)
}

func LoadDBSpecificEndpoints(router *routing.RouteGroup) {
	router.Post("/service-groups/register", authentication.AdminRequiredMiddleware, controllers.RegisterServiceGroupHandler)
	router.Post("/service-groups/<group_id>/services", authentication.AdminRequiredMiddleware, controllers.AddServiceToGroupHandler)
	// sd.sdAPI.Post("/service-groups/service/register", authentication.AuthorizationMiddleware, RegisterServiceToServiceGroupHandler)
	router.Get("/service-groups", controllers.ListServiceGroupsHandler)
}
