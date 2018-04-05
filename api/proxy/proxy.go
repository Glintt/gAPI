package proxy

import (
	"gAPIManagement/api/authentication"
	"gAPIManagement/api/cache"
	"fmt"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery"

	"github.com/qiangxue/fasthttp-routing"
)

var sd servicediscovery.ServiceDiscovery
var oauthserver authentication.OAuthServer

func StartProxy(router *routing.Router) {
	oauthserver = authentication.LoadFromConfig()

	router.To("GET,POST,PUT,PATCH,DELETE", "/*", HandleRequest)

	sd.SetIsService(false)
}

func HandleRequest(c *routing.Context) error {
	fmt.Println("\n=============================================================")
	fmt.Println("\n====================== NEW REQUEST ==========================")
	fmt.Println("\n=============================================================")
	fmt.Println("PROXY =====> Method = " + string(c.Method()) + "; URI = " + string(c.Request.RequestURI()))
	fmt.Println("===============================================================")

	cachedRequest := cache.GetCacheForRequest(c)

	if cachedRequest.Service.ToURI == "" {
		fmt.Println("SD NOT FROM CACHE")

		var err error
		cachedRequest.Service, err = getServiceFromServiceDiscovery(c)

		if err != nil {
			http.Response(c, `{"error": true, "msg": "Resource not found"}`, 404, "/proxy")
			return nil
		}

		cachedRequest.UpdateServiceCache = true
	} else {
		fmt.Println("SD FROM CACHE")
	}

	if cachedRequest.Protection.UserInfo == "" {
		fmt.Println("PROTECTION NOT FROM CACHE")
		cachedRequest.Protection = checkAuthorization(c, cachedRequest.Service)

		if cachedRequest.Protection.Error != nil {
			http.Response(c, `{"error":true, "msg":"Not Authorized."}`, 401, cachedRequest.Service.MatchingURI)
			return nil
		}

		cachedRequest.UpdateProtectionCache = true
	} else {
		fmt.Println("PROTECTION FROM CACHE")
	}

	if cachedRequest.Response.StatusCode == 0 {
		fmt.Println("RESPONSE NOT FROM CACHE")
		cachedRequest.Response = getApiResponse(c, cachedRequest.Protection, cachedRequest.Service)
		if cachedRequest.Response.StatusCode < 300 {
			cachedRequest.UpdateResponseCache = true
		}
	} else {
		fmt.Println("RESPONSE FROM CACHE")
	}


	http.Response(c, string(cachedRequest.Response.Body), cachedRequest.Response.StatusCode, cachedRequest.Service.MatchingURI)

	cache.StoreRequestInfoToCache(c, cachedRequest)
	return nil
}

func getApiResponse(c *routing.Context, authorization authentication.ProtectionInfo, service servicediscovery.Service) http.ResponseInfo {

	c.Request.Header.Set(authorization.Header, authorization.UserInfo)
	headers := http.GetHeadersFromRequest(c.Request)

	response := service.Call(string(c.Method()), string(c.Request.RequestURI()), headers, "")

	return http.ResponseInfo{
		StatusCode:  response.Header.StatusCode(),
		ContentType: response.Header.ContentType(),
		Body:        response.Body()}
}

func checkAuthorization(c *routing.Context, service servicediscovery.Service) authentication.ProtectionInfo {
	if !service.Protected {
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
