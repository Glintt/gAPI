package http

import (
	"gAPIManagement/api/utils"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func GetHeadersFromRequest(request fasthttp.Request) map[string]string {
	var headersMap map[string]string
	headersMap = make(map[string]string)

	request.Header.VisitAll(func(key []byte, val []byte) {
		headersMap[string(key)] = string(val)
	})

	return headersMap
}

func GetQueryParamsFromRequest(c *routing.Context) map[string]string {
	var queryParamsMap map[string]string
	queryParamsMap = make(map[string]string)

	c.QueryArgs().VisitAll(func(key []byte, val []byte) {
		queryParamsMap[string(key)] = string(val)
	})

	return queryParamsMap
}

func GetQueryParamsFromRequestCtx(c *fasthttp.RequestCtx) map[string]string {
	var queryParamsMap map[string]string
	queryParamsMap = make(map[string]string)

	c.QueryArgs().VisitAll(func(key []byte, val []byte) {
		queryParamsMap[string(key)] = string(val)
	})

	return queryParamsMap
}

func GetURIWithParams(c *routing.Context) string {
	queryParams := GetQueryParamsFromRequest(c)

	uri := string(c.Request.RequestURI())

	uri = strings.Trim(uri, " ")
	uri = strings.TrimSuffix(uri, "%20")
	uri = strings.TrimPrefix(uri, "%20")
	uri = strings.Trim(uri, " ")

	for pKey, pValue := range queryParams {
		if pValue == "" {
			uri = strings.Replace(uri, pKey+"=", pKey, 1)
			uri = strings.Replace(uri, pKey, pKey+"=", 1)
		}
	}

	return uri
}

func addHeadersToRequest(request *fasthttp.Request, headers map[string]string) {
	if headers == nil {
		return
	}
	for key, value := range headers {
		if key != "N/a" {
			request.Header.Set(key, value)
		}
	}
}

func addQueryParamsToRequest(request *fasthttp.Request, queryParams map[string]string) {
	if queryParams == nil {
		return
	}
}

func MakeRequest(method string, url string, body string, headers map[string]string) *fasthttp.Response {
	request := fasthttp.AcquireRequest()
	request.SetRequestURI(url)
	request.Header.SetMethod(method)

	//request.Header.SetContentType("application/json")
	request.SetBody([]byte(body))

	addHeadersToRequest(request, headers)
	if _, ok := headers["Content-Type"]; ok {
		request.Header.SetContentType(headers["Content-Type"])
	}
	utils.LogMessage("=============================================================", utils.DebugLogType)
	utils.LogMessage("HTTP Request ---- Method: "+method+" ; Url = "+url+" ; Body = "+body, utils.DebugLogType)
	utils.LogMessage("             ---- Headers: "+request.Header.String(), utils.DebugLogType)
	utils.LogMessage("=============================================================", utils.DebugLogType)

	client := fasthttp.Client{}

	resp := fasthttp.AcquireResponse()
	err := client.Do(request, resp)

	if err != nil {
		utils.LogMessage(err.Error(), utils.ErrorLogType)
		resp.SetStatusCode(400)
	}

	fasthttp.ReleaseRequest(request)

	return resp
}
