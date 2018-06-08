package routes

import (
	"github.com/qiangxue/fasthttp-routing"
	"gAPIManagement/api/controllers"
)

func InitAuthenticationAPI(router *routing.Router) {
	router.Post("/oauth/token", controllers.GetTokenHandler)
	router.Get("/oauth/authorize", controllers.AuthorizeTokenHandler)
	router.Get("/oauth/me", controllers.MeHandler)
}