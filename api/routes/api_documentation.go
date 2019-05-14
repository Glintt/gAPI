package routes

import (
	"gAPIManagement/api/controllers"

	routing "github.com/qiangxue/fasthttp-routing"
)

func InitServiceDocumentationRoutes(router *routing.Router) {
	router.To("GET,POST,PUT,PATCH,DELETE", "/api_docs/<service_name>/documentation", controllers.HandleServiceDocumentationRequest)
	router.To("GET,POST,PUT,PATCH,DELETE", "/api_docs/<service_name>/*", controllers.HandleServiceDocumentationJSRequest)
}
