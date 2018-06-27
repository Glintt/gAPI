package proxy

import (
	"strconv"
	"gAPIManagement/api/ratelimiting"
	"gAPIManagement/api/utils"
	"regexp"
	"strings"
	"gAPIManagement/api/cache"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery"
	authentication "gAPIManagement/api/thirdpartyauthentication"

	"github.com/qiangxue/fasthttp-routing"
)

var sd servicediscovery.ServiceDiscovery
var oauthserver authentication.OAuthServer

var SERVICE_NAME = "/proxy"

func StartProxy(router *routing.Router) {
	oauthserver = authentication.LoadFromConfig()

	ratelimiting.InitRateLimiting()
	router.To("GET,POST,PUT,PATCH,DELETE", "/*", ratelimiting.RateLimiting, HandleRequest)

	sd = *servicediscovery.GetServiceDiscoveryObject()
}

func HandleRequest(c *routing.Context) error {
	utils.LogMessage("=========================================", utils.InfoLogType)
	utils.LogMessage("REQUEST =====> Method = "+string(c.Method())+"; URI = "+string(c.Request.RequestURI()), utils.InfoLogType)
	utils.LogMessage("=========================================", utils.InfoLogType)

	cachedRequest := cache.GetCacheForRequest(c)

	if cachedRequest.Service.ToURI == "" {
		utils.LogMessage("SD NOT FROM CACHE", utils.DebugLogType)

		var err error
		cachedRequest.Service, err = getServiceFromServiceDiscovery(c)

		utils.LogMessage("IsExternalRequest = " + strconv.FormatBool(sd.IsExternalRequest(c)), utils.DebugLogType)
		utils.LogMessage("IsReachableFromExternal = " + strconv.FormatBool(cachedRequest.Service.IsReachableFromExternal(sd)), utils.DebugLogType)

		if err != nil || (sd.IsExternalRequest(c) && !cachedRequest.Service.IsReachableFromExternal(sd)) {
			http.Response(c, `{"error": true, "msg": "Resource not found"}`, 404, SERVICE_NAME)
			return nil
		}

		cachedRequest.UpdateServiceCache = true
	} else {
		utils.LogMessage("SD FROM CACHE", utils.DebugLogType)
	}

	if !cachedRequest.Protection.Cached {
		utils.LogMessage("PROTECTION NOT FROM CACHE", utils.DebugLogType)
		cachedRequest.Protection = checkAuthorization(c, cachedRequest.Service)

		if cachedRequest.Protection.Error != nil {
			http.Response(c, `{"error":true, "msg":"Not Authorized."}`, 401, cachedRequest.Service.MatchingURI)
			return nil
		}

		cachedRequest.UpdateProtectionCache = true
	} else {
		utils.LogMessage("PROTECTION FROM CACHE", utils.DebugLogType)
	}

	if cachedRequest.Response.StatusCode == 0 {
		utils.LogMessage("RESPONSE NOT FROM CACHE", utils.DebugLogType)
		cachedRequest.Response = getApiResponse(c, cachedRequest.Protection, cachedRequest.Service)
		if cachedRequest.Response.StatusCode < 300 {
			cachedRequest.UpdateResponseCache = true
		}
	} else {
		utils.LogMessage("RESPONSE FROM CACHE", utils.DebugLogType)
	}

	http.Response(c, string(cachedRequest.Response.Body), cachedRequest.Response.StatusCode, cachedRequest.Service.MatchingURI)

	if cachedRequest.Response.StatusCode < 300 {
		cache.StoreRequestInfoToCache(c, cachedRequest)
	}
	return nil
}

func getApiResponse(c *routing.Context, authorization authentication.ProtectionInfo, service servicediscovery.Service) http.ResponseInfo {

	c.Request.Header.Set(authorization.Header, authorization.UserInfo)
	headers := http.GetHeadersFromRequest(c.Request)
	if _, ok := headers["X-Forwarded-For"]; !ok {
		headers["X-Forwarded-For"] = c.RemoteIP().String()
	}

	body := c.Request.Body()
	
	response := service.Call(string(c.Method()), http.GetURIWithParams(c) , headers, string(body))

	return http.ResponseInfo{
		StatusCode:  response.Header.StatusCode(),
		ContentType: response.Header.ContentType(),
		Body:        response.Body()}
}

func checkAuthorization(c *routing.Context, service servicediscovery.Service) authentication.ProtectionInfo {
	if !service.Protected {
		return authentication.ProtectionInfo{Header: "n/a", UserInfo: "n/a", Error: nil}
	}

	endpoint := strings.Replace(string(c.RequestURI()), service.MatchingURI, "", 1)
	endpoint = strings.Split(endpoint, "?")[0]

	if _, ok := service.ProtectedExclude[endpoint]; !ok {
		for key, val := range service.ProtectedExclude {
			re := regexp.MustCompile(key)

			value := strings.ToLower(val)
			method := strings.ToLower(string(c.Method()))

			if re.ReplaceAllString(endpoint, "") == "" && strings.Contains(value, method) {
				return authentication.ProtectionInfo{Header: "n/a", UserInfo: "n/a", Error: nil}
			}
		}
	} else {
		return authentication.ProtectionInfo{Header: "n/a", UserInfo: "n/a", Error: nil}
	}

	headerName, userInfo, notAuthorizedErr := Protect(service, c)

	if notAuthorizedErr != nil {
		return authentication.ProtectionInfo{Header: "n/a", UserInfo: "n/a", Error: notAuthorizedErr}
	}

	return authentication.ProtectionInfo{Header: headerName, UserInfo: userInfo, Error: nil}

}

func getServiceFromServiceDiscovery(c *routing.Context) (servicediscovery.Service, error) {
	return sd.GetEndpointForUri(string(c.Request.RequestURI()))
}

func Protect(service servicediscovery.Service, c *routing.Context) (string, string, error) {
	if service.Protected {
		return oauthserver.Authorize(c.Request)
	}
	return "n/a", "n/a", nil
}
