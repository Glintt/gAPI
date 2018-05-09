package servicediscovery

import (
	"fmt"
	"gAPIManagement/api/http"
	"strings"
	"regexp"
	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

type Service struct {
	ID                    bson.ObjectId `bson:"_id" json:"id"`
	Name                  string
	Domain                string
	Port                  string
	MatchingURI           string
	ToURI                 string
	Protected             bool
	APIDocumentation      string
	CachingExpirationTime int
	IsCachingActive       bool
	IsActive              bool
	HealthcheckUrl	string
	LastActiveTime int64
	RestartHost       string
	RestartPort       string
	RestartEndpoint       string
}

func (service *Service) Call(method string, uri string, headers map[string]string, body string) *fasthttp.Response {
	uri = strings.Replace(uri, service.MatchingURI, service.ToURI, 1)

	callURLWithoutProtocol := service.Domain + ":" + service.Port + uri
	callURLWithoutProtocol = strings.Replace(callURLWithoutProtocol, "//", "/", -1)

	callURL := "http://" + callURLWithoutProtocol

	return http.MakeRequest(method, callURL, body, headers)
}

func (service *Service) Restart() bool {
	method := "POST"
	callURL := "http://" + service.RestartHost + ":" + service.RestartPort + service.RestartEndpoint

	var validURL = regexp.MustCompile(`^(((http|https):\/{2})+(([0-9a-z_-]+\.?)+(:[0-9]+)?((\/([~0-9a-zA-Z#\+%@\.\/_-]+))?(\?[0-9a-zA-Z\+%@\/&\[\];=_-]+)?)?))\b$`)

	if validURL.MatchString(callURL) {
		resp := http.MakeRequest(method, callURL, "", nil)
		fmt.Println(string(resp.Body()))
		if resp.StatusCode() != 200 {
			return false
		}
		return true
	}
	return false
}