package main

import (
	"encoding/json"
	"runtime"
	"strconv"

	apianalytics "github.com/Glintt/gAPI/api/api-analytics"
	"github.com/Glintt/gAPI/api/routes"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/users"

	"os"
	"time"

	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/cache"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/healthcheck"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/logs"
	"github.com/Glintt/gAPI/api/plugins"
	"github.com/Glintt/gAPI/api/proxy"
	"github.com/Glintt/gAPI/api/servicediscovery"
	"github.com/Glintt/gAPI/api/sockets"
	"github.com/Glintt/gAPI/api/utils"
	"github.com/joho/godotenv"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var server = &fasthttp.Server{}
var router *routing.Router

var defaultHttpsPort = "443"

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.LogMessage("Error loading .env file", utils.ErrorLogType)
	}

	if os.Getenv("GO_MAX_PROCS") != "" {
		maxProcs, err := strconv.Atoi(os.Getenv("GO_MAX_PROCS"))
		if err != nil {
			runtime.GOMAXPROCS(maxProcs)
		}
	}

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

	utils.LogMessage("Listening on port: "+listeningPort, utils.InfoLogType)
	utils.LogMessage("Using HTTPS: "+strconv.FormatBool(config.GApiConfiguration.Protocol.Https), utils.InfoLogType)

	if config.GApiConfiguration.Protocol.Https {
		go func() {
			panic(fasthttp.ListenAndServeTLS(":"+defaultHttpsPort, config.GApiConfiguration.Protocol.CertificateFile, config.GApiConfiguration.Protocol.CertificateKey, CORSHandle))
		}()
	}

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

	if serviceString == constants.SERVICE_NAME || serviceString == proxy.SERVICE_NAME || serviceString == apianalytics.SERVICE_NAME || serviceString == authentication.SERVICE_NAME || serviceString == users.SERVICE_NAME || serviceString == plugins.SERVICE_NAME {
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
		indexName = config.GAPI_API_LOGS_INDEX
	}

	utils.LogMessage("Log IndexName = "+indexName, utils.DebugLogType)

	elapsedTime := utils.CurrentTimeMilliseconds() - beginTime
	queryArgs, _ := json.Marshal(http.GetQueryParamsFromRequestCtx(ctx))
	headersObject := http.GetHeadersFromRequest(ctx.Request)
	delete(headersObject, "Authorization")
	headers, _ := json.Marshal(headersObject)
	logRequest := logs.NewRequestLogging(ctx, queryArgs, headers, utils.CurrentDateWithFormat(time.UnixDate), elapsedTime, string(service), indexName)
	work := logs.LogWorkRequest{Name: "", LogToSave: logRequest}
	logs.WorkQueue <- work
}
