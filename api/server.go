package main

import (
	apianalytics "gAPIManagement/api/api-analytics"
	"gAPIManagement/api/cache"
	"gAPIManagement/api/config"
	"gAPIManagement/api/logs"
	"gAPIManagement/api/utils"
	"gAPIManagement/api/proxy"
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/authentication"
	"gAPIManagement/api/http"
	"fmt"
	"os"
	"time"
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var server = &fasthttp.Server{}
var router *routing.Router

func main() {
	config.LoadConfigs()

	router = routing.New()
	router.Get("/reload",authentication.AuthorizationMiddleware, ReloadServices)
	router.Get("/invalidate-cache",authentication.AuthorizationMiddleware, InvalidateCache)
	initServices()

	listenAPI(router)
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

func initServices() {
	authentication.InitGAPIAuthenticationServer(router)
	cache.InitCachingService()
	logs.StartDispatcher(2)
	servicediscovery.StartServiceDiscovery(router)
	apianalytics.StartAPIAnalytics(router)
	proxy.StartProxy(router)
}

func listenAPI(router *routing.Router) {
	listeningPort := os.Getenv("API_MANAGEMENT_PORT")

	if listeningPort == "" {
		listeningPort = "8080"
	}

	fmt.Println("Listening on port: " + listeningPort)
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

	LogRequest(ctx, beginTime)
}

func LogRequest(ctx *fasthttp.RequestCtx, beginTime int64){
	if ! config.GApiConfiguration.Logs.Active || string(ctx.Method()) == "OPTIONS" {
		return
	}

	elapsedTime := utils.CurrentTimeMilliseconds() - beginTime
	service := ctx.Response.Header.Peek("service")
	queryArgs, _ := json.Marshal(http.GetQueryParamsFromRequestCtx(ctx))
	headers, _ := json.Marshal(http.GetHeadersFromRequest(ctx.Request))
	logRequest := logs.NewRequestLogging(ctx, queryArgs, headers, utils.CurrentDateWithFormat(time.UnixDate), elapsedTime, string(service))
	work := logs.LogWorkRequest{Name: "", LogToSave: logRequest}
	logs.WorkQueue <- work
}