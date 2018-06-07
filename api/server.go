package main

import (
	"gAPIManagement/api/users"
	"encoding/json"
	
	apianalytics "gAPIManagement/api/api-analytics"
	"gAPIManagement/api/authentication"
	"gAPIManagement/api/cache"
	"gAPIManagement/api/config"
	"gAPIManagement/api/healthcheck"
	"gAPIManagement/api/http"
	"gAPIManagement/api/logs"
	"gAPIManagement/api/proxy"
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/sockets"
	"gAPIManagement/api/utils"
	"os"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var server = &fasthttp.Server{}
var router *routing.Router

func main() {
	config.LoadConfigs()
	
	router = routing.New()
	router.Get("/reload", authentication.AuthorizationMiddleware, ReloadServices)
	router.Get("/invalidate-cache", authentication.AuthorizationMiddleware, InvalidateCache)
	initServices()

	InitSocketServices()
	
	listenAPI(router)
}

func InitSocketServices() {
	go sockets.SocketListen()
	sockets.StartRequestsCounterSender()
}

func InvalidateCache(c *routing.Context) error {
	cache.InvalidateCache()
	c.Response.SetBody([]byte(`{"error":false, "msg": "Invalidation finished."}`))
	c.Response.Header.SetContentType("application/json")
	return nil
}

func ReloadServices(c *routing.Context) error {
	initServices()
	cache.InvalidateCache()
	c.Response.SetBody([]byte(`{"error":false, "msg": "Reloaded successfully."}`))
	c.Response.Header.SetContentType("application/json")
	return nil
}

func InitUsersService(router *routing.Router) {
	users.InitUsers()
	
	usersGroup := router.Group(config.USERS_GROUP)

	usersGroup.Get("", authentication.AdminRequiredMiddleware, users.FindUsersHandler)
	usersGroup.Get("/", authentication.AdminRequiredMiddleware, users.FindUsersHandler)
	usersGroup.Post("/", authentication.AdminRequiredMiddleware, users.CreateUserHandler)
	usersGroup.Post("", authentication.AdminRequiredMiddleware, users.CreateUserHandler)
	usersGroup.Get("/<username>", authentication.AdminRequiredMiddleware, users.GetUserHandler)
}

func initServices() {
	authentication.InitGAPIAuthenticationServer(router)
	cache.InitCachingService()
	logs.StartDispatcher(2)
	servicediscovery.StartServiceDiscovery(router)
	InitUsersService(router)
	apianalytics.StartAPIAnalytics(router)
	proxy.StartProxy(router)
	healthcheck.InitHealthCheck()
}

func listenAPI(router *routing.Router) {
	listeningPort := os.Getenv("API_MANAGEMENT_PORT")

	if listeningPort == "" {
		listeningPort = "8080"
	}

	utils.LogMessage("Listening on port: " + listeningPort)
	panic(fasthttp.ListenAndServe(":"+listeningPort, CORSHandle))
}

var (
	corsAllowHeaders     = "access-control-allow-origin, Content-Type, Authorization"
	corsAllowMethods     = "HEAD, GET, POST, PUT, DELETE, OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func CORSHandle(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
	ctx.Response.Header.Set("Content-Type", "application/json")

	beginTime := utils.CurrentTimeMilliseconds()
	router.HandleRequest(ctx)
	service := ctx.Response.Header.Peek("service")

	RequestCounterSocket(service)
	LogRequest(ctx, service, beginTime)
}

func RequestCounterSocket(service []byte) {
	sockets.IncrementRequestCounter()
}

func LogRequest(ctx *fasthttp.RequestCtx, service []byte, beginTime int64) {
	if !config.GApiConfiguration.Logs.Active || string(ctx.Method()) == "OPTIONS" {
		return
	}

	elapsedTime := utils.CurrentTimeMilliseconds() - beginTime
	queryArgs, _ := json.Marshal(http.GetQueryParamsFromRequestCtx(ctx))
	headers, _ := json.Marshal(http.GetHeadersFromRequest(ctx.Request))
	logRequest := logs.NewRequestLogging(ctx, queryArgs, headers, utils.CurrentDateWithFormat(time.UnixDate), elapsedTime, string(service))
	work := logs.LogWorkRequest{Name: "", LogToSave: logRequest}
	logs.WorkQueue <- work
}
