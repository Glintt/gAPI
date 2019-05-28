package routes

import (
	"github.com/Glintt/gAPI/api/controllers"
	"github.com/qiangxue/fasthttp-routing"
)


func InitAppsGroupsApi(router *routing.Router) {
	router.Get("/apps-groups", controllers.GetAppGroups)
	router.Post("/apps-groups", controllers.CreateAppGroup)
	router.Get("/apps-groups/matches", controllers.AppGroupsMatches)
	router.Get("/apps-groups/ungrouped", controllers.UngroupedApps)
	router.Get("/apps-groups/search/<service_id>", controllers.FindAppGroupForService)
	router.Put("/apps-groups/<group_id>", controllers.UpdateAppGroup)
	router.Get("/apps-groups/<group_id>", controllers.GetAppGroupById)	
	router.Delete("/apps-groups/<group_id>", controllers.DeleteAppGroup)
	router.Post("/apps-groups/<group_id>/<service_id>", controllers.AssociateServiceToAppGroup)
	router.Delete("/apps-groups/<group_id>/<service_id>", controllers.DeassociateServiceFromApplicationGroup)
}