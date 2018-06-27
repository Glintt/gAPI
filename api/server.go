package main

import (
	"gAPIManagement/api/users"
	"gAPIManagement/api/api-analytics"
	"gAPIManagement/api/routes"
	"encoding/json"
	
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
	
	InitServices()
	
	InitAPIs()

	InitSocketServices()
	
	listenAPI(router)
}

func InitSocketServices() {
	go sockets.SocketListen()
	sockets.StartRequestsCounterSender()
}


func InitAPIs() {
	routes.InitAPIRoutes(router)

	proxy.StartProxy(router)	
}

func InitServices() {
	cache.InitCachingService()
	if config.GApiConfiguration.Logs.Active {
		logs.StartDispatcher(2)
	}
	authentication.InitGAPIAuthenticationServer()
	servicediscovery.InitServiceDiscovery()
	healthcheck.InitHealthCheck()
}

func listenAPI(router *routing.Router) {
	listeningPort := os.Getenv("API_MANAGEMENT_PORT")

	if listeningPort == "" {
		listeningPort = "8080"
	}

	utils.LogMessage("Listening on port: " + listeningPort, utils.InfoLogType)
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
	
	if string(ctx.Request.Header.Peek("Connection")) != "keep-alive" {
		defer ctx.Response.SetConnectionClose()
	}
}

func RequestCounterSocket(service []byte) {
	sockets.IncrementRequestCounter()
}

func IsGApiService(service []byte) bool {
	serviceString := string(service)
	
	if serviceString == servicediscovery.SERVICE_NAME || serviceString == proxy.SERVICE_NAME || serviceString == apianalytics.SERVICE_NAME || serviceString == authentication.SERVICE_NAME || serviceString == users.SERVICE_NAME {
		return true
	}

	return false
}

func LogRequest(ctx *fasthttp.RequestCtx, service []byte, beginTime int64) {
	if !config.GApiConfiguration.Logs.Active || string(ctx.Method()) == "OPTIONS" {
		return
	}

	indexName := ""
	if IsGApiService(service) {
		indexName = "gapi-api-logs"
	}
	
	elapsedTime := utils.CurrentTimeMilliseconds() - beginTime
	queryArgs, _ := json.Marshal(http.GetQueryParamsFromRequestCtx(ctx))
	headers, _ := json.Marshal(http.GetHeadersFromRequest(ctx.Request))
	logRequest := logs.NewRequestLogging(ctx, queryArgs, headers, utils.CurrentDateWithFormat(time.UnixDate), elapsedTime, string(service), indexName)
	work := logs.LogWorkRequest{Name: "", LogToSave: logRequest}
	logs.WorkQueue <- work
}
