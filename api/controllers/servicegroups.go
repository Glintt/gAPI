package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery"
	"github.com/Glintt/gAPI/api/servicediscovery/servicegroup"
	"github.com/Glintt/gAPI/api/authentication"

	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

func ListServiceGroupsHandler(c *routing.Context) error {
	user := authentication.GetAuthenticatedUser(c)
	sg, err := ServiceDiscovery(user).GetListOfServicesGroup()

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	json, _ := json.Marshal(sg)
	http.Response(c, string(json), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func RegisterServiceGroupHandler(c *routing.Context) error {

	serviceGroup, err := servicediscovery.ValidateServiceGroupBody(c)
	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	serviceGroup.Id = bson.NewObjectId()

	err = servicediscovery.ServiceGroupMethods()["create"].(func(servicegroup.ServiceGroup) error)(serviceGroup)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service created successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func AddServiceToGroupHandler(c *routing.Context) error {
	serviceGroupId := c.Param("group_id")

	var bodyMap map[string]string
	err := json.Unmarshal(c.Request.Body(), &bodyMap)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	if _, ok := bodyMap["service_id"]; !ok {
		http.Response(c, `{"error": "Invalid body. Missing service_id."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	if serviceGroupId == "null" || bodyMap["service_id"] == "null" || bodyMap["service_id"] == "" {
		http.Response(c, `{"error": "Invalid body."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	err = servicediscovery.ServiceGroupMethods()["addservicetogroup"].(func(string, string) error)(serviceGroupId, bodyMap["service_id"])

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service added to group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func DeassociateServiceFromGroup(c *routing.Context) error {
	serviceGroupId := c.Param("group_id")
	serviceId := c.Param("service_id")

	if serviceGroupId == "null" || serviceId == "null" || serviceId == "" {
		http.Response(c, `{"error": "Invalid body."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	err := servicediscovery.ServiceGroupMethods()["removeservicefromgroup"].(func(string, string) error)(serviceGroupId, serviceId)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service deassociated from group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateServiceGroup(c *routing.Context) error {
	serviceGroupId := c.Param("group_id")

	var sGroup servicegroup.ServiceGroup
	sgNew := c.Request.Body()
	json.Unmarshal(sgNew, &sGroup)

	err := servicediscovery.ServiceGroupMethods()["update"].(func(string, servicegroup.ServiceGroup) error)(
		serviceGroupId, sGroup)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service group update successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func RemoveServiceGroup(c *routing.Context) error {
	serviceGroupId := c.Param("group_id")

	err := servicediscovery.ServiceGroupMethods()["delete"].(func(string) error)(serviceGroupId)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Service group removed successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func GetServiceGroupHandler(c *routing.Context) error {
	serviceGroup := string(c.Param("group"))

	sg, err := servicediscovery.ServiceGroupMethods()["getbyid"].(func(string) (servicegroup.ServiceGroup, error))(serviceGroup)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	sgByte, _ := json.Marshal(sg)
	http.Response(c, string(sgByte), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}
