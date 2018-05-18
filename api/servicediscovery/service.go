package servicediscovery

import (
	"gopkg.in/mgo.v2/bson"
	"gAPIManagement/api/http"
	"strings"

	"github.com/valyala/fasthttp"
	//"gopkg.in/mgo.v2/bson"
)

type Service struct {
	Id                    bson.ObjectId `bson:"_id" json:"Id"`
	Name                  string
	Domain                string
	Port                  string
	MatchingURI           string
	MatchingURIRegex      string
	ToURI                 string
	Protected             bool
	APIDocumentation      string
	CachingExpirationTime int
	IsCachingActive       bool
	IsActive              bool
	HealthcheckUrl        string
	LastActiveTime        int64
}

func (service *Service) Call(method string, uri string, headers map[string]string, body string) *fasthttp.Response {
	uri = strings.Replace(uri, service.MatchingURI, service.ToURI, 1)

	callURLWithoutProtocol := service.Domain + ":" + service.Port + uri
	callURLWithoutProtocol = strings.Replace(callURLWithoutProtocol, "//", "/", -1)

	callURL := "http://" + callURLWithoutProtocol

	return http.MakeRequest(method, callURL, body, headers)
}

func (service *Service) GenerateId() bson.ObjectId {
	return bson.NewObjectId()
}

func (service *Service) NormalizeService() {
	service.MatchingURIRegex = GetMatchingURIRegex(service.MatchingURI)
	if service.Id == "" {
		service.Id = service.GenerateId()
	}
}