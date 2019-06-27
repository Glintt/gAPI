package controllers

import (
	"io/ioutil"
	"net/http"
	"crypto/tls"
	netUrl "net/url"
	"path"
	"strings"

	gapiHttp "github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	"github.com/Glintt/gAPI/api/utils"

	routing "github.com/qiangxue/fasthttp-routing"
)

func HandleServiceDocumentationRequest(c *routing.Context) error {
	url, err := GetServiceDocumentationUrl(c.Param("service_name"))
	if err != nil {
		gapiHttp.Response(c, "Not found", 404, "", "text/html")
		return nil
	}

	response, responseString, err := GetHtml(url)
	if err != nil {
		gapiHttp.Response(c, responseString, 500, "", "text/html")
		return nil
	}
	gapiHttp.Response(c, responseString, 200, "", response.Header.Get("Content-Type"))

	return nil
}

func HandleServiceDocumentationJSRequest(c *routing.Context) error {
	serviceIdentifier := c.Param("service_name")
	url, err := GetServiceDocumentationUrl(serviceIdentifier)

	if err != nil {
		gapiHttp.Response(c, "Not found", 404, "", "text/html")
		return nil
	}

	uri := strings.Replace(string(c.Request.RequestURI()), "/api_docs/"+serviceIdentifier+"/", "", -1)
	u, _ := netUrl.Parse(url)

	uri = path.Join(u.Path, uri)
	url = u.Scheme + "://" + u.Host + uri

	response, responseString, err := GetHtml(url)
	if err != nil {
		gapiHttp.Response(c, responseString, 500, "", "text/html")
		return nil
	}

	gapiHttp.Response(c, responseString, 200, "", response.Header.Get("Content-Type"))

	return nil
}

func GetServiceDocumentationUrl(serviceIdentifier string) (string, error) {
	service, err := servicediscovery.GetInternalServiceDiscoveryObject().FindService(service.Service{Identifier: serviceIdentifier})
	if err != nil {
		return "", err
	}

	return service.APIDocumentation, nil
}

func GetHtml(url string) (*http.Response, string, error) {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
	client := &http.Client{Transport: tr}
	
	response, err := client.Get(url)
	if err != nil {
		utils.LogMessage(err.Error(), utils.ErrorLogType)
		return response, err.Error(), err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.LogMessage(err.Error(), utils.ErrorLogType)
		return response, err.Error(), err
	}

	return response, string(responseData), nil
}
